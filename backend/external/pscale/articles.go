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
func (aDB *ArticlesDB) Get(offset, limit int) (*core.ArticleSeries, error) {
	series := make(core.ArticleSeries, 0, limit)

	query := "SELECT * FROM articles ORDER BY published_on DESC LIMIT ? OFFSET ?"

	rows, err := aDB.DB.Queryx(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("PScaleArticles: unable to query articles table - %w", err)
	}

	for rows.Next() {
		var article core.Article
		var questions, questionDisplay, topics string
		err = rows.Scan(&article.ID, &article.Title, &article.URL, &topics, &questions, &questionDisplay, &article.PublishedOn)
		if err != nil {
			return nil, fmt.Errorf("PScaleArticles: error scanning row - %w", err)
		}
		article.Questions = strings.Split(questions, "\n")
		article.QuestionDisplay = strings.Split(questionDisplay, "\n")
		article.Topics = strings.Split(topics, ",")

		article.Date = time.Unix(article.PublishedOn, 0).Format("Jan 2, 2006")

		series = append(series, article)
	}

	return &series, nil
}

// Find implements a mysql fulltext search of the articles table
func (aDB *ArticlesDB) Find(terms string) (*core.ArticleSeries, error) {
	series := make(core.ArticleSeries, 0, 12)

	query := "SELECT * FROM articles WHERE MATCH (title, topics, question_display) AGAINST (? IN BOOLEAN MODE) ORDER BY id DESC"

	rows, err := aDB.DB.Queryx(query, terms)
	if err != nil {
		return nil, fmt.Errorf("PScaleArticles: unable to query pscale database - %w", err)
	}

	for rows.Next() {
		var article core.Article
		var questions, questionDisplay, topics string
		err = rows.Scan(&article.ID, &article.Title, &article.URL, &topics, &questions, &questionDisplay, &article.PublishedOn)
		if err != nil {
			return nil, fmt.Errorf("PScaleArticles: error scanning row - %w", err)
		}
		article.Questions = strings.Split(questions, "\n")
		article.QuestionDisplay = strings.Split(questionDisplay, "\n")
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
	query := "INSERT INTO articles (title, url, topics, questions, question_display, published_on) VALUES (?, ?, ?, ?, ?, ?)"

	for _, article := range *data {
		tx, err := aDB.DB.Begin()
		if err != nil {
			return fmt.Errorf("PScaleArticles: unable to begin tx for adding articles input to db - %w", err)
		}
		defer tx.Rollback()

		res, err := tx.Exec(query, article.Title, article.URL, strings.Join(article.Topics, ","), strings.Join(article.Questions, "\n"), strings.Join(article.QuestionDisplay, "\n"), article.PublishedOn)
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

		rx := regexp.MustCompile(`\d{4}\s-\sQ\d{1,2}`)
		for _, x := range article.Questions {
			question := rx.FindString(x)

			_, err = tx.Exec("INSERT INTO questions (question, article_id) VALUES (?, ?)", question, id)
			if err != nil {
				return fmt.Errorf("PScaleArticles: unable to add questions for article %s to db - %w", article.Title, err)
			}
		}

		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("PScaleArticles: unable to commit tx to insert articles to db - %w", err)
		}
	}

	return nil
}

func (aDB *ArticlesDB) Update(data *core.ArticleSeries) error {
	query := "UPDATE articles SET title = ?, url = ?, topics = ?, questions = ?, question_display= ?, published_on = ? WHERE id = ?"

	for _, article := range *data {
		tx, err := aDB.DB.Begin()
		if err != nil {
			return fmt.Errorf("PScaleArticles: unable to begin tx for updating articles in db - %w", err)
		}
		defer tx.Rollback()

		res, err := tx.Exec(query, article.Title, article.URL, strings.Join(article.Topics, ","), strings.Join(article.Questions, "\n"), strings.Join(article.QuestionDisplay, "\n"), article.PublishedOn, article.ID)
		if err != nil {
			return fmt.Errorf("PScaleArticles: unable to update article %s in db - %w", article.Title, err)
		}

		id, _ := res.LastInsertId()
		for _, w := range article.Topics {
			_, err = tx.Exec("INSERT INTO topics (topic, article_id) VALUES (?, ?)", w, id)
			if err != nil {
				return fmt.Errorf("PScaleArticles: unable to update topics for article %s in db - %w", article.Title, err)
			}
		}

		rx := regexp.MustCompile(`\d{4}\s-\sQ\d{1,2}`)
		for _, x := range article.Questions {
			question := rx.FindString(x)
			_, err = tx.Exec("INSERT INTO questions (question, article_id) VALUES (?, ?)", question, id)
			if err != nil {
				return fmt.Errorf("PScaleArticles: unable to update questions for article %s in db - %w", article.Title, err)
			}
		}

		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("PScaleArticles: unable to commit tx to update articles in db - %w", err)
		}
	}

	return nil
}

func (aDB *ArticlesDB) Delete(ids string) error {
	query := fmt.Sprintf("DELETE FROM articles WHERE id in (%s)", ids)

	tx, err := aDB.DB.Begin()
	if err != nil {
		return fmt.Errorf("PScaleArticles: unable to begin tx for deleting articles from db - %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(query)
	if err != nil {
		return fmt.Errorf("PScaleArticles: unable to delete articles from db - %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("PScaleArticles: unable to commit tx to delet articles from db - %w", err)
	}

	return nil
}

func (aDB *ArticlesDB) GetQuestion(qn string) (string, error) {
	query := fmt.Sprintf("SELECT wording FROM question_list WHERE question = ?")

	row := aDB.DB.QueryRowx(query, qn)

	var wording string
	err := row.Scan(&wording)
	if err != nil {
		return "", fmt.Errorf("error scanning row to get user- %w", err)
	}

	question := fmt.Sprintf("%s %s", qn, wording)

	return question, nil
}
