package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jwnpoh/njcreaderapp/cockroach-migration-tool/core"
)

func migrateArticlesTable(pscaleDB, cockroachDB *sqlx.DB) error {
	fmt.Println("getting articles from pscale...")
	articles, err := getPscaleArticles(pscaleDB)
	if err != nil {
		return fmt.Errorf("PScaleArticles: unable to get articles from pscale - %w", err)
	}
	fmt.Printf("got %d articles from pscale\n", len(articles))

	fmt.Printf("attempting to insert %d articles from pscale to cockcroach...\n", len(articles))
	err = insertArticlesToCockroach(cockroachDB, articles)
	if err != nil {
		return fmt.Errorf("CockroachArticles: unable to insert articles to cockroach - %w\n", err)
	}

	return nil
}

func getPscaleArticles(pscaleDB *sqlx.DB) (core.ArticleSeries, error) {
	series := make(core.ArticleSeries, 0)

	query := "SELECT * FROM articles"

	fmt.Println("running query")
	rows, err := pscaleDB.Queryx(query)
	if err != nil {
		return nil, fmt.Errorf("PScaleArticles: unable to query articles table - %w\n", err)
	}

	var count int
	for rows.Next() {
		var article core.Article
		var id int
		var questions, questionDisplay, topics string
		err = rows.Scan(&id, &article.Title, &article.URL, &topics, &questions, &questionDisplay, &article.PublishedOn, &article.MustRead)
		if err != nil {
			return nil, fmt.Errorf("PScaleArticles: error scanning row - %w\n", err)
		}
		article.Questions = strings.Split(questions, "\n")
		article.QuestionDisplay = strings.Split(questionDisplay, "\n")
		article.Topics = strings.Split(topics, ",")

		article.Date = time.Unix(article.PublishedOn, 0).Format("Jan 2, 2006")

		series = append(series, article)
		count++
		fmt.Printf("scanned %d articles...\n", count)
	}

	fmt.Printf("scanned a total of %d articles from pscale\n", len(series))

	return series, nil
}

func insertArticlesToCockroach(cockroachDB *sqlx.DB, articles core.ArticleSeries) error {
	if len(articles) < 1 {
		fmt.Println("did not receive articles to insert to cockroach.")
		return fmt.Errorf("did not receive articles to insert to cockroach.")
	}

	query := "INSERT INTO articles (title, url, topics, questions, question_display, published_on, must_read) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"

	fmt.Println("running query")
	tx, err := cockroachDB.Begin()
	if err != nil {
		return fmt.Errorf("Cockroach: unable to begin tx for adding articles input to cockroach - %w\n", err)
	}
	defer tx.Rollback()

	for i, article := range articles {
		fmt.Printf("adding article #%d\n", i+1)
		var id string
		err := tx.QueryRow(query, article.Title, article.URL, strings.Join(article.Topics, ","), strings.Join(article.Questions, "\n"), strings.Join(article.QuestionDisplay, "\n"), article.PublishedOn, article.MustRead).Scan(&id)
		if err != nil {
			return fmt.Errorf("CockroachArticles: unable to add article %s to db - %v\n", article.Title, err)
		}
		fmt.Printf("string form of uuid is %s\n", id)
		uuid, _ := uuid.Parse(id)

		fmt.Println("inserting into topics table...")
		for _, w := range article.Topics {
			_, err = tx.Exec("INSERT INTO topics (topic, article_id) VALUES ($1, $2)", w, uuid)
			if err != nil {
				return fmt.Errorf("PScaleArticles: unable to add topics for article %s to db - %w\n", article.Title, err)
			}
		}

		rx := regexp.MustCompile(`\d{4}\s-\sQ\d{1,2}`)
		for _, x := range article.Questions {
			question := rx.FindString(x)

			_, err = tx.Exec("INSERT INTO questions (question, article_id) VALUES ($1, $2)", question, uuid)
			if err != nil {
				return fmt.Errorf("PScaleArticles: unable to add questions for article %s to db - %w", article.Title, err)
			}
		}

		fmt.Printf("successfully inserted article #%d with uuid %v\n", i+1, uuid)
		if i == len(articles)-1 {
			err = tx.Commit()
			if err != nil {
				return fmt.Errorf("PScaleArticles: unable to commit tx to insert articles to db - %w", err)
			}
		}
	}

	return nil
}
