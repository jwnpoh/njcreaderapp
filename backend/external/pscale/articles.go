package pscale

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
)

type ArticlesDB struct {
	DB *sqlx.DB
}

// NewPscaleDB returns a connection interface for the application to connect to the planetscale database.
func NewArticlesDB(dsn string) (*ArticlesDB, error) {
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("PScaleArticles: unable to initialize pscale database - %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("PScaleArticles: no response from pscale database - %w", err)
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)

	return &ArticlesDB{DB: db}, nil
}

// Get retrieves a slice of 10 articles from the planetscale database with limit and offset in the query.
func (aDB *ArticlesDB) Get(offset int) (*core.ArticleSeries, error) {
	series := make(core.ArticleSeries, 0, 10)

	query := "SELECT * FROM articles ORDER BY id DESC LIMIT 10 OFFSET ?"

	rows, err := aDB.DB.Queryx(query, offset)
	if err != nil {
		return nil, fmt.Errorf("PScaleArticles: unable to query articles table for page %d - %w", offset, err)
	}

	for rows.Next() {
		var article core.Article
		var questions, topics string
		err = rows.Scan(&article.ID, &article.Title, &article.URL, &topics, &questions, &article.PublishedOn)
		if err != nil {
			return nil, fmt.Errorf("PScaleArticles: error scanning row - %w", err)
		}
		// need to fix. some questions contain a comma in the question wording and should not be split there.
		article.Questions = strings.Split(questions, "\n")
		article.Topics = strings.Split(topics, ",")

		article.Date = time.Unix(article.PublishedOn, 0).Format("Jan 2, 2006")

		series = append(series, article)
	}

	return &series, nil
}

// Find implements a mysql fulltext search of the articles table
func (aDB *ArticlesDB) Find(terms string) (*core.ArticleSeries, error) {
	series := make(core.ArticleSeries, 0, 12)

	query := "SELECT * FROM articles WHERE MATCH (title, topics, questions) AGAINST (? IN BOOLEAN MODE) ORDER BY id DESC"

	rows, err := aDB.DB.Queryx(query, terms)
	if err != nil {
		return nil, fmt.Errorf("PScaleArticles: unable to query pscale database - %w", err)
	}

	for rows.Next() {
		var article core.Article
		var questions, topics string
		err = rows.Scan(&article.ID, &article.Title, &article.URL, &topics, &questions, &article.PublishedOn)
		if err != nil {
			return nil, fmt.Errorf("PScaleArticles: error scanning row - %w", err)
		}
		// need to fix. some questions contain a comma in the question wording and should not be split there.
		article.Questions = strings.Split(questions, "\n")
		article.Topics = strings.Split(topics, ",")

		article.Date = time.Unix(article.PublishedOn, 0).Format("Jan 2, 2006")

		series = append(series, article)
	}

	if len(series) == 0 {
		return &series, fmt.Errorf("PScaleArticles: no articles matched the query %s", terms)
	}

	return &series, nil
}

// Store stores a slice of articles sent from the front end admin dashboard via the articles service.
func (aDB *ArticlesDB) Store(data *core.ArticleSeries) error {
	query := "INSERT INTO articles (title, url, topics, questions, published_on) VALUES (?, ?, ?, ?, ?)"

	for _, article := range *data {
		tx, err := aDB.DB.Begin()
		if err != nil {
			return fmt.Errorf("PScaleArticles: unable to begin tx for adding articles input to db - %w", err)
		}
		defer tx.Rollback()

		res, err := tx.Exec(query, article.Title, article.URL, strings.Join(article.Topics, ","), strings.Join(article.Questions, ","), article.PublishedOn)
		if err != nil {
			return fmt.Errorf("PScaleArticles: unable to add article %s to db - %w", article.Title, err)
		}

		id, _ := res.LastInsertId()
		for _, w := range article.Topics {
			_, err = tx.Exec("INSERT INTO topics (topic, article_id) VALUES (?, ?)", w, id)
			if err != nil {
				return fmt.Errorf("PScaleArticles: unable to add topics for article %s to db - %w", article.Title, err)
			}
		}

		rx := regexp.MustCompile(`\d{1,4}\s\d{1,2}`)
		for _, x := range article.Questions {
			question := rx.Find([]byte(x))

			_, err = tx.Exec("INSERT INTO questions (question, article_id) VALUES (?, ?)", question, id)
			if err != nil {
				return fmt.Errorf("PScaleArticles: unable to add questions for article %s to db - %w", article.Title, err)
			}
		}

		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("PScaleArticles: unable to commit tx to db - %w", err)
		}
	}

	return nil
}
