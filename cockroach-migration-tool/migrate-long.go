package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/jwnpoh/njcreaderapp/cockroach-migration-tool/core"
)

func migrateLongTable(pscaleDB, cockroachDB *sqlx.DB) error {
	fmt.Println("getting long articles from pscale...")
	articles, err := getPscaleLong(pscaleDB)
	if err != nil {
		return fmt.Errorf("PScaleLong: unable to get long articles from pscale - %w", err)
	}
	fmt.Printf("got %d articles from pscale\n", len(articles))

	fmt.Printf("attempting to insert %d articles from pscale to cockcroach...\n", len(articles))
	err = insertLongToCockroach(cockroachDB, articles)
	if err != nil {
		return fmt.Errorf("CockroachArticles: unable to insert articles to cockroach - %w\n", err)
	}

	return nil
}

func getPscaleLong(pscaleDB *sqlx.DB) (core.LongSeries, error) {
	series := make(core.LongSeries, 0)

	query := "SELECT * FROM long"

	fmt.Println("running query")
	rows, err := pscaleDB.Queryx(query)
	if err != nil {
		return nil, fmt.Errorf("PScaleLong: unable to query long table - %w\n", err)
	}

	var count int
	for rows.Next() {
		var long core.Long
		var id int
		err = rows.Scan(&id, &long.Title, &long.URL, &long.Topic)
		if err != nil {
			return nil, fmt.Errorf("PScaleLong: error scanning row - %w\n", err)
		}

		series = append(series, long)
		count++
		fmt.Printf("scanned %d long articles...\n", count)
	}

	fmt.Printf("scanned a total of %d articles from pscale\n", len(series))

	return series, nil
}

func insertLongToCockroach(cockroachDB *sqlx.DB, longs core.LongSeries) error {
	if len(longs) < 1 {
		fmt.Println("did not receive longs to insert to cockroach.")
		return fmt.Errorf("did not receive longs to insert to cockroach.")
	}

	query := "INSERT INTO long (title, url, topic) VALUES ($1, $2, $3)"

	fmt.Println("running query")
	tx, err := cockroachDB.Begin()
	if err != nil {
		return fmt.Errorf("Cockroach: unable to begin tx for adding articles input to cockroach - %w\n", err)
	}
	defer tx.Rollback()

	for i, long := range longs {
		fmt.Printf("adding long article #%d\n", i+1)
		_, err := tx.Exec(query, long.Title, long.URL, long.Topic)
		if err != nil {
			return fmt.Errorf("CockroachLong: unable to add long article %s to db - %v\n", long.Title, err)
		}

		fmt.Printf("successfully inserted article #%d\n", i+1)
		if i == len(longs)-1 {
			err = tx.Commit()
			if err != nil {
				return fmt.Errorf("PScaleArticles: unable to commit tx to insert articles to db - %w", err)
			}
		}
	}

	return nil
}
