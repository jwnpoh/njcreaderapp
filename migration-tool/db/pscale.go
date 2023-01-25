package db

import (
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func newDB(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return db, fmt.Errorf("unable to initialize sql database - %w", err)
	}

	err = db.Ping()
	if err != nil {
		return db, fmt.Errorf("no response from sql database - %w", err)
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)

	return db, nil
}

func PscaleAdd(a *ArticlesDBByDate, db *sqlx.DB) error {
	defer db.Close()

	var count int

	tx, err := db.Begin()
	if err != nil {
		// return fmt.Errorf("unable to begin tx for adding article %s to db - %w", article.Title, err)
		return fmt.Errorf("unable to begin tx for adding articles to db - %w", err)
	}
	defer tx.Rollback()

	for i, article := range *a {

		// var questions []string
		var questions, questionDisplay string
		q := strings.Builder{}
		t := strings.Builder{}
		for i, v := range article.Questions {
			q.WriteString(v.Year + " - Q" + v.Number)
			t.WriteString(v.Year + " - Q" + v.Number + " " + v.Wording)
			if i < len(article.Questions)-1 {
				q.WriteString("\n")
				t.WriteString("\n")
			}
			questions = q.String()
			questionDisplay = t.String()
		}

		res, err := tx.Exec("INSERT INTO articles (title, url, topics, questions, question_display, published_on, must_read) VALUES (?, ?, ?, ?, ?, ?, ?)", article.Title, article.URL, strings.Join(article.Topics, ","), questions, questionDisplay, article.Date, false)
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
			_, err = tx.Exec("INSERT INTO questions (question, article_id) VALUES (?, ?)", fmt.Sprintf("%v - Q%v", x.Year, x.Number), id)
			if err != nil {
				return fmt.Errorf("unable to add questions for article %s to db - %w", article.Title, err)
			}
		}

		if i >= 0 && i%200 == 0 || i == len(*a)-1 {
			err = tx.Commit()
			if err != nil {
				return fmt.Errorf("unable to commit article %s to db - %w", article.Title, err)
			}
			tx, _ = db.Begin()
			defer tx.Rollback()
		}
		count++
		fmt.Printf("Added %v/%v articles\r", count, a.Len())
	}
	return nil
}

func PscaleAddQuestions(qnDB QuestionsDB, db *sqlx.DB) error {
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("unable to begin tx for adding questions to db - %w", err)
	}
	defer tx.Rollback()

	for _, v := range qnDB {
		_, err := tx.Exec("INSERT INTO question_list (question, year, number, wording) VALUES (?, ?, ?, ?)", fmt.Sprintf("%s - Q%s", v.Year, v.Number), v.Year, v.Number, v.Wording)
		if err != nil {
			return fmt.Errorf("unable to add question %v to db - %w", v, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("unable to commit questions to db - %w", err)
	}

	return nil
}
