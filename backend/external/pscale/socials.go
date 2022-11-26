package pscale

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type SocialsDB struct {
	DB *sqlx.DB
}

// NewPscaleDB returns a connection interface for the application to connect to the planetscale database.
func NewSocialsDB(dsn string) (*SocialsDB, error) {
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("PScalePosts: unable to initialize pscale database - %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("PScalePosts: no response from pscale database - %w", err)
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)

	return &SocialsDB{DB: db}, nil
}

func (sDB *SocialsDB) GetFollowing(userID int) ([]int, error) {
	following := make([]int, 0)

	query := "SELECT follows FROM following WHERE user_id = ?"

	rows, err := sDB.DB.Queryx(query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return following, nil
		}
		return nil, fmt.Errorf("PScaleSocials: unable to query following table - %w", err)
	}

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("PScaleSocials: error scanning row - %w", err)
		}

		following = append(following, id)
	}

	return following, nil
}

func (sDB *SocialsDB) GetFollowedBy(userID int) ([]int, error) {
	followedBy := make([]int, 0)

	query := "SELECT follows FROM followed_by WHERE user_id = ?"

	rows, err := sDB.DB.Queryx(query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return followedBy, nil
		}
		return nil, fmt.Errorf("PScaleSocials: unable to query followed_by table - %w", err)
	}

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("PScaleSocials: error scanning row - %w", err)
		}

		followedBy = append(followedBy, id)
	}

	return followedBy, nil
}
