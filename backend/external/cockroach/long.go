package cockroach

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
)

type LongsDB struct {
	DB *sqlx.DB
}

// NewPscaleDB returns a connection interface for the application to connect to the planetscale database.
func NewLongsDB(dsn string) (*LongsDB, error) {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("PScaleLong: unable to initialize pscale database - %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("PScaleLong: no response from pscale database - %w", err)
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)

	return &LongsDB{DB: db}, nil
}

func (lDB *LongsDB) GetTopics() (*core.LongTopics, error) {
	topics := make(core.LongTopics, 0)

	query := "SELECT DISTINCT topic FROM long ORDER BY topic"

	rows, err := lDB.DB.Queryx(query)
	if err != nil {
		return nil, fmt.Errorf("PScaleLong: unable to query pscale database - %w", err)
	}

	for rows.Next() {
		var topic string
		err = rows.Scan(&topic)
		if err != nil {
			return nil, fmt.Errorf("PScaleLong: error scanning row - %w", err)
		}
		topics = append(topics, topic)
	}
	return &topics, nil
}

func (lDB *LongsDB) Get(topic string) (*core.LongSeries, error) {
	if topic == "all" {
		return lDB.GetAll()
	}

	longs := make(core.LongSeries, 0)

	query := "SELECT * FROM long WHERE topic = $1"

	if topic == "all" {
		query = "SELECT * FROM long"
	}

	rows, err := lDB.DB.Queryx(query, topic)
	if err != nil {
		return nil, fmt.Errorf("PScaleLong: unable to query pscale database - %w", err)
	}

	for rows.Next() {
		var long core.Long
		err = rows.Scan(&long.ID, &long.Title, &long.URL, &long.Topic)
		if err != nil {
			return nil, fmt.Errorf("PScaleLong: error scanning row - %w", err)
		}
		longs = append(longs, long)
	}

	return &longs, nil
}

func (lDB *LongsDB) GetAll() (*core.LongSeries, error) {
	longs := make(core.LongSeries, 0)

	query := "SELECT * FROM long"

	rows, err := lDB.DB.Queryx(query)
	if err != nil {
		return nil, fmt.Errorf("PScaleLong: unable to query pscale database - %w", err)
	}

	for rows.Next() {
		var long core.Long
		err = rows.Scan(&long.ID, &long.Title, &long.URL, &long.Topic)
		if err != nil {
			return nil, fmt.Errorf("PScaleLong: error scanning row - %w", err)
		}
		longs = append(longs, long)
	}

	return &longs, nil
}
func (lDB *LongsDB) Store(data *core.LongPayload) error {
	query := "INSERT INTO long (title, url, topic) VALUES ($1, $2, $3)"

	tx, err := lDB.DB.Begin()
	if err != nil {
		return fmt.Errorf("PScaleLong: unable to begin tx to add long articles to db - %w", err)
	}
	defer tx.Rollback()

	for i, long := range *data {
		_, err = tx.Exec(query, long.Title, long.URL, long.Topic)
		if err != nil {
			return fmt.Errorf("PScaleLong: unable to add long article %s to db - %w", long.Title, err)
		}

		if i >= 0 && i%200 == 0 || i == len(*data)-1 {
			err = tx.Commit()
			if err != nil {
				return fmt.Errorf("PScaleLong: unable to commit tx to store long articlesn db - %w", err)
			}
			tx, err = lDB.DB.Begin()
			if err != nil {
				return fmt.Errorf("PScaleLong: unable to begin tx to add long articles to db - %w", err)
			}
			defer tx.Rollback()
		}
	}
	return nil
}

func (lDB *LongsDB) Update(data *core.Long) error {
	query := "UPDATE long SET title = $1, url = $2, topic = $3 WHERE id = $4"

	tx, err := lDB.DB.Begin()
	if err != nil {
		return fmt.Errorf("PScaleLong: unable to begin tx to add long articles to db - %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(query, data.Title, data.URL, data.Topic, data.ID)
	if err != nil {
		return fmt.Errorf("PScaleLong: unable to add long article %s to db - %w", data.Title, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("PScaleLong: unable to commit tx to update long article %s - %w", data.Title, err)
	}

	return nil
}

func (lDB *LongsDB) Delete(ids string) error {
	query := fmt.Sprintf("DELETE FROM long WHERE id in (%s)", ids)

	tx, err := lDB.DB.Begin()
	if err != nil {
		return fmt.Errorf("PScaleLong: unable to begin tx for deleting long articles from db - %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(query)
	if err != nil {
		return fmt.Errorf("PScaleLong: unable to delete long articles from db - %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("PScaleLong: unable to commit tx to delete long articles from db - %w", err)
	}

	return nil
}
