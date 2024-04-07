package cockroach

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
)

type ArticlesDB struct {
	DB *sqlx.DB
}

// NewArticlesDB returns a connection interface for the application to connect to the cockroachDB database.
func NewArticlesDB(dsn string) (*ArticlesDB, error) {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("CockroachArticles: unable to initialize cockroach database - %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("cockroachArticles: no response from cockroach database - %w", err)
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)

	return &ArticlesDB{DB: db}, nil
}

// Get retrieves a slice of 10 articles from the planetscale database with limit and offset in the query.
func (aDB *ArticlesDB) Get(offset, limit int) (*core.ArticleSeries, error) {
	series := make(core.ArticleSeries, 0, limit)

	query := "SELECT * FROM articles ORDER BY published_on DESC LIMIT $1 OFFSET $2"

	rows, err := aDB.DB.Queryx(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("cockroachArticles: unable to query articles table - %w", err)
	}

	for rows.Next() {
		var article core.Article
		var questions, questionDisplay, topics string
		err = rows.Scan(&article.ID, &article.Title, &article.URL, &topics, &questions, &questionDisplay, &article.PublishedOn, &article.MustRead)
		if err != nil {
			return nil, fmt.Errorf("cockroachArticles: error scanning row - %w", err)
		}
		article.Questions = strings.Split(questions, "\n")
		article.QuestionDisplay = strings.Split(questionDisplay, "\n")
		article.Topics = strings.Split(topics, ",")

		article.Date = time.Unix(article.PublishedOn, 0).Format("Jan 2, 2006")

		series = append(series, article)
	}

	return &series, nil
}

func (aDB *ArticlesDB) GetArticle(id int) (*core.Article, error) {
	query := "SELECT * FROM articles WHERE id = $1"

	row := aDB.DB.QueryRowx(query, id)

	var article core.Article
	var questions, questionDisplay, topics string
	err := row.Scan(&article.ID, &article.Title, &article.URL, &topics, &questions, &questionDisplay, &article.PublishedOn, &article.MustRead)
	if err != nil {
		return nil, fmt.Errorf("cockroachArticles: error scanning row - %w", err)
	}
	article.Questions = strings.Split(questions, "\n")
	article.QuestionDisplay = strings.Split(questionDisplay, "\n")
	article.Topics = strings.Split(topics, ",")

	article.Date = time.Unix(article.PublishedOn, 0).Format("Jan 2, 2006")

	return &article, nil
}

// Find implements a mysql fulltext search of the articles table
func (aDB *ArticlesDB) Find(terms string) (*core.ArticleSeries, error) {
	series := make(core.ArticleSeries, 0, 12)

	// TODO:
	// NEED TO UPDATE FUNCTION!!
	query := "SELECT * FROM articles WHERE MATCH (title, topics, question_display) AGAINST ($1 IN BOOLEAN MODE) ORDER BY id DESC"

	rows, err := aDB.DB.Queryx(query, terms)
	if err != nil {
		return nil, fmt.Errorf("cockroachArticles: unable to query cockroach database - %w", err)
	}

	for rows.Next() {
		var article core.Article
		var questions, questionDisplay, topics string
		err = rows.Scan(&article.ID, &article.Title, &article.URL, &topics, &questions, &questionDisplay, &article.PublishedOn, &article.MustRead)
		if err != nil {
			return nil, fmt.Errorf("cockroachArticles: error scanning row - %w", err)
		}
		article.Questions = strings.Split(questions, "\n")
		article.QuestionDisplay = strings.Split(questionDisplay, "\n")
		article.Topics = strings.Split(topics, ",")

		article.Date = time.Unix(article.PublishedOn, 0).Format("Jan 2, 2006")

		series = append(series, article)
	}

	if len(series) == 0 {
		return &series, fmt.Errorf("cockroachArticles: no articles matched the query %s", terms)
	}

	return &series, nil
}

// Store stores a slice of articles sent from the front end admin dashboard via the articles service.
func (aDB *ArticlesDB) Store(data *core.ArticleSeries) error {
	// query := "INSERT INTO articles (title, url, topics, questions, question_display, published_on, must_read) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	query := "INSERT INTO articles (title, url, topics, questions, question_display, published_on, must_read) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"

	tx, err := aDB.DB.Begin()
	if err != nil {
		return fmt.Errorf("cockroachArticles: unable to begin tx for adding articles input to db - %w", err)
	}
	defer tx.Rollback()

	for i, article := range *data {
		var id string
		err := tx.QueryRow(query, article.Title, article.URL, strings.Join(article.Topics, ","), strings.Join(article.Questions, "\n"), strings.Join(article.QuestionDisplay, "\n"), article.PublishedOn, article.MustRead).Scan(&id)
		if err != nil {
			return fmt.Errorf("CockroachArticles: unable to add article %s to db - %v\n", article.Title, err)
		}
		// res, err := tx.Exec(query, article.Title, article.URL, strings.Join(article.Topics, ","), strings.Join(article.Questions, "\n"), strings.Join(article.QuestionDisplay, "\n"), article.PublishedOn, article.MustRead)
		// if err != nil {
		// 	return fmt.Errorf("cockroachArticles: unable to add article %s to db - %w", article.Title, err)
		// }
		uuid, _ := uuid.Parse(id)

		for _, w := range article.Topics {
			_, err = tx.Exec("INSERT INTO topics (topic, article_id) VALUES ($1, $2)", w, uuid)
			if err != nil {
				return fmt.Errorf("cockroachArticles: unable to add topics for article %s to db - %w", article.Title, err)
			}
		}

		rx := regexp.MustCompile(`\d{4}\s-\sQ\d{1,2}`)
		for _, x := range article.Questions {
			question := rx.FindString(x)

			_, err = tx.Exec("INSERT INTO questions (question, article_id) VALUES ($1, $2)", question, uuid)
			if err != nil {
				return fmt.Errorf("cockroachArticles: unable to add questions for article %s to db - %w", article.Title, err)
			}
		}

		if i == len(*data)-1 {
			err = tx.Commit()
			if err != nil {
				return fmt.Errorf("cockroachArticles: unable to commit tx to insert articles to db - %w", err)
			}
		}
	}

	return nil
}

func (aDB *ArticlesDB) Update(data *core.ArticleSeries) error {
	query := "UPDATE articles SET title = $1, url = $2, topics = $3, questions = $4, question_display= $5, published_on = $6, must_read = $7 WHERE id = $8"

	tx, err := aDB.DB.Begin()
	if err != nil {
		return fmt.Errorf("cockroachArticles: unable to begin tx for updating articles in db - %w", err)
	}
	defer tx.Rollback()

	for i, article := range *data {
		// res, err := tx.Exec(query, article.Title, article.URL, strings.Join(article.Topics, ","), strings.Join(article.Questions, "\n"), strings.Join(article.QuestionDisplay, "\n"), article.PublishedOn, article.MustRead, article.ID)
		var id string
		uuid, _ := uuid.Parse(id)
		fmt.Printf("updating article uuid %s\ntitle %s\ntopics %s\nquestions %s\n", uuid, article.Title, strings.Join(article.Topics, ","), strings.Join(article.Questions, "\n"))
		_, err := tx.Exec(query, article.Title, article.URL, strings.Join(article.Topics, ","), strings.Join(article.Questions, "\n"), strings.Join(article.QuestionDisplay, "\n"), article.PublishedOn, article.MustRead, uuid)
		if err != nil {
			return fmt.Errorf("CockroachArticles: unable to add article %s to db - %v\n", article.Title, err)
		}
		if err != nil {
			return fmt.Errorf("cockroachArticles: unable to update article %s in db - %w", article.Title, err)
		}

		for _, w := range article.Topics {
			_, err = tx.Exec("INSERT INTO topics (topic, article_id) VALUES ($1, $2)", w, uuid)
			if err != nil {
				return fmt.Errorf("cockroachArticles: unable to update topics for article %s in db - %w", article.Title, err)
			}
		}

		rx := regexp.MustCompile(`\d{4}\s-\sQ\d{1,2}`)
		for _, x := range article.Questions {
			question := rx.FindString(x)
			_, err = tx.Exec("INSERT INTO questions (question, article_id) VALUES ($1, $2)", question, uuid)
			if err != nil {
				return fmt.Errorf("cockroachArticles: unable to update questions for article %s in db - %w", article.Title, err)
			}
		}

		if i == len(*data)-1 {
			err = tx.Commit()
			if err != nil {
				return fmt.Errorf("cockroachArticles: unable to commit tx to update articles in db - %w", err)
			}
		}
	}

	return nil
}

func (aDB *ArticlesDB) Delete(ids []string) error {
	var uuids []uuid.UUID

	query := "DELETE FROM articles WHERE id = $1"

	for _, id := range ids {
		uuid, err := uuid.Parse(id)
		if err != nil {
			return fmt.Errorf("cockroachArticles: unable to parse uuid from string uuid - %w", err)
		}
		uuids = append(uuids, uuid)

		tx, err := aDB.DB.Begin()
		if err != nil {
			return fmt.Errorf("cockroachArticles: unable to begin tx for deleting articles from db - %w", err)
		}
		defer tx.Rollback()

		_, err = tx.Exec(query, uuid)
		if err != nil {
			return fmt.Errorf("cockroachArticles: unable to delete articles from db - %w", err)
		}

		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("cockroachArticles: unable to commit tx to delet articles from db - %w", err)
		}
	}

	return nil
}

func (aDB *ArticlesDB) GetQuestion(qn string) (string, error) {
	query := fmt.Sprintf("SELECT wording FROM question_list WHERE question = $1")

	row := aDB.DB.QueryRowx(query, qn)

	var wording string
	err := row.Scan(&wording)
	if err != nil {
		return "", fmt.Errorf("error scanning row to get user- %w", err)
	}

	question := fmt.Sprintf("%s %s", qn, wording)

	return question, nil
}
