package pscale

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
)

type StatsDB struct {
	DB *sqlx.DB
}

// NewPscaleDB returns a connection interface for the application to connect to the planetscale database.
func NewStatsDB(dsn string) (*StatsDB, error) {
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

	return &StatsDB{DB: db}, nil
}

func (sDB *StatsDB) GetStats() (*core.Stats, error) {
	res := core.Stats{}

	// get number of articles
	query := "SELECT COUNT(*) FROM articles"

	row := sDB.DB.QueryRowx(query)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("PScaleArticles: error scanning row - %w", err)
	}

	res.NumberofArticles = count

	// get topics with most articles
	query2 := "SELECT topic, COUNT(*) FROM topics GROUP BY topic ORDER BY COUNT(*) DESC LIMIT 5"

	rows, err := sDB.DB.Queryx(query2)
	if err != nil {
		return nil, fmt.Errorf("PScaleArticles: unable to query topics table - %w", err)
	}

	for rows.Next() {
		var topic string
		var count int
		var kv core.KV
		err = rows.Scan(&topic, &count)
		if err != nil {
			return nil, fmt.Errorf("PScaleArticles: error scanning row for topic stats - %w", err)
		}
		kv.K = topic
		kv.V = count
		res.TopicsWithMostArticles = append(res.TopicsWithMostArticles, kv)
	}

	// get topics with fewest articles
	query3 := "SELECT topic, COUNT(*) FROM topics GROUP BY topic ORDER BY COUNT(*) LIMIT 5"

	rows, err = sDB.DB.Queryx(query3)
	if err != nil {
		return nil, fmt.Errorf("PScaleArticles: unable to query topics table - %w", err)
	}

	for rows.Next() {
		var topic string
		var count int
		var kv core.KV
		err = rows.Scan(&topic, &count)
		if err != nil {
			return nil, fmt.Errorf("PScaleArticles: error scanning row for topic stats - %w", err)
		}
		kv.K = topic
		kv.V = count
		res.TopicsWithFewestArticles = append(res.TopicsWithFewestArticles, kv)
	}

	// get questions with most articles
	query4 := "SELECT question, COUNT(*) FROM questions GROUP BY question ORDER BY COUNT(*) DESC LIMIT 5"

	rows, err = sDB.DB.Queryx(query4)
	if err != nil {
		return nil, fmt.Errorf("PScaleArticles: unable to query questions table - %w", err)
	}

	for rows.Next() {
		var question string
		var count int
		var kv core.KV
		err = rows.Scan(&question, &count)
		if err != nil {
			return nil, fmt.Errorf("PScaleArticles: error scanning row for question stats - %w", err)
		}
		kv.K = question
		kv.V = count
		res.QuestionsWithMostArticles = append(res.QuestionsWithMostArticles, kv)
	}

	// get questions with fewest articles
	query5 := "SELECT question, COUNT(*) FROM questions GROUP BY question ORDER BY COUNT(*) LIMIT 5"

	rows, err = sDB.DB.Queryx(query5)
	if err != nil {
		return nil, fmt.Errorf("PScaleArticles: unable to query questions table - %w", err)
	}

	for rows.Next() {
		var question string
		var count int
		var kv core.KV
		err = rows.Scan(&question, &count)
		if err != nil {
			return nil, fmt.Errorf("PScaleArticles: error scanning row for question stats - %w", err)
		}
		kv.K = question
		kv.V = count
		res.QuestionsWthFewestArticles = append(res.QuestionsWthFewestArticles, kv)
	}

	return &res, nil
}
