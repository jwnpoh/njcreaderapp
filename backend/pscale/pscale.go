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

type PscaleDB struct {
	DB *sqlx.DB
}

func NewPscaleDB() (*PscaleDB, error) {
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

	return &PscaleDB{DB: db}, nil
}

func (pScale *PscaleDB) Get(page int) (*core.ArticleSeries, error) {
	series := make(core.ArticleSeries, 0, 12)

	query := "SELECT * FROM articles ORDER BY id DESC LIMIT 12 OFFSET ?"

	rows, err := pScale.DB.Queryx(query, page)
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
