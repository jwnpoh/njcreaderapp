package pscale

import (
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
)

// PScale provides interface for services to connect to the planetscale database.
type PScale interface {
	Get(page int) (*core.ArticleSeries, error)
	Find(term string) (*core.ArticleSeries, error)
	Store(data *core.ArticleSeries) error
	GetQns(questions []string) (*[]core.Question, error)
}

type pscaleDB struct {
	DB *sqlx.DB
}

// NewPscaleDB returns a connection interface for the application to connect to the planetscale database.
func NewPscaleDB() (PScale, error) {
	db, err := sqlx.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize pscale database - %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("no response from pscale database - %w", err)
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)

	return &pscaleDB{DB: db}, nil
}

// Get retrieves a slice of 12 articles from the planetscale database with limit and offset in the query.
func (ps *pscaleDB) Get(offset int) (*core.ArticleSeries, error) {
	series := make(core.ArticleSeries, 0, 12)

	query := "SELECT * FROM articles ORDER BY id DESC LIMIT 12 OFFSET ?"

	rows, err := ps.DB.Queryx(query, offset)
	if err != nil {
		return nil, fmt.Errorf("unable to query articles table for page %d - %w", offset, err)
	}

	for rows.Next() {
		var article core.Article
		var questions, topics string
		err = rows.Scan(&article.ID, &article.Title, &article.URL, &topics, &questions, &article.PublishedOn)
		if err != nil {
			return nil, fmt.Errorf("error scanning row - %w", err)
		}
		article.Questions = strings.Split(questions, ",")
		article.Topics = strings.Split(topics, ",")

		article.Date = time.Unix(article.PublishedOn, 0).Format("Jan 2, 2006")

		series = append(series, article)
	}

	return &series, nil
}

// Find implements a search of the articles table and the question_list table for a match of the search term
func (ps *pscaleDB) Find(term string) (*core.ArticleSeries, error) {
	term = "%" + term + "%"
	series := make(core.ArticleSeries, 0, 12)

	// Match questions with term and find article id
	query := "SELECT question FROM question_list WHERE wording LIKE ?"
	var b strings.Builder
	rows, err := ps.DB.Queryx(query, term)
	if err != nil {
		return nil, fmt.Errorf("unable to query question_list table for the search term '%s' - %w", term, err)
	}
	for rows.Next() {
		var q string
		err = rows.Scan(&q)
		if err != nil {
			return nil, fmt.Errorf("error scanning query results into Go variable - %w", err)
		}
		fmt.Fprintf(&b, "%s|", q)
	}

	// Match articles with term in title or topics

	query = "SELECT * FROM articles WHERE title LIKE ? OR topics LIKE ? OR questions REGEXP ?"

	rows, err = ps.DB.Queryx(query, term, term, b.String())
	if err != nil {
		return nil, fmt.Errorf("unable to query pscale database - %w", err)
	}

	for rows.Next() {
		var article core.Article
		var questions, topics string
		err = rows.Scan(&article.ID, &article.Title, &article.URL, &topics, &questions, &article.PublishedOn)
		if err != nil {
			return nil, fmt.Errorf("error scanning row - %w", err)
		}
		article.Questions = strings.Split(questions, ",")
		article.Topics = strings.Split(topics, ",")

		article.Date = time.Unix(article.PublishedOn, 0).Format("Jan 2, 2006")

		series = append(series, article)
	}

	return &series, nil
}

func (ps *pscaleDB) GetQns(questions []string) (*[]core.Question, error) {

	qns := make([]core.Question, 0)

	query := "SELECT year, number, wording FROM question_list WHERE (question = ?)"

	for _, q := range questions {
		var qn core.Question
		err := ps.DB.Get(&qn, query, q)
		if err != nil {
			return nil, fmt.Errorf("error scanning row to retrieve question %s - %w", q, err)
		}
		qns = append(qns, qn)
	}
	return &qns, nil
}

// Store stores a slice of articles sent from the front end admin dashboard via the articles service.
func (ps *pscaleDB) Store(data *core.ArticleSeries) error {

	for _, article := range *data {
		tx, err := ps.DB.Begin()

		if err != nil {
			return fmt.Errorf("unable to begin tx for adding articles input to db - %w", err)
		}
		defer tx.Rollback()

		res, err := tx.Exec("INSERT INTO articles (title, url, topics, questions, published_on) VALUES (?, ?, ?, ?, ?)", article.Title, article.URL, strings.Join(article.Topics, ","), strings.Join(article.Questions, ","), article.PublishedOn)
		if err != nil {
			return fmt.Errorf("unable to add article %s to db - %w", article.Title, err)
		}

		id, _ := res.LastInsertId()
		for _, w := range article.Topics {
			_, err = tx.Exec("INSERT INTO topics (topic, article_id) VALUES (?, ?)", w, id)
			if err != nil {
				return fmt.Errorf("unable to add topics for article %s to db - %w", article.Title, err)
			}
		}

		for _, x := range article.Questions {
			_, err = tx.Exec("INSERT INTO questions (question, article_id) VALUES (?, ?)", x, id)
			if err != nil {
				return fmt.Errorf("unable to add questions for article %s to db - %w", article.Title, err)
			}
		}

		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("unable to commit tx to db - %w", err)
		}
	}

	return nil
}
