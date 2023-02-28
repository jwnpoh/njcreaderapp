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

	query := "SELECT follows FROM follows WHERE user_id = ?"

	rows, err := sDB.DB.Queryx(query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return following, nil
		}
		return nil, fmt.Errorf("PScaleSocials: unable to get follows - %w", err)
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

	query := "SELECT user_id FROM follows WHERE follows = ?"

	rows, err := sDB.DB.Queryx(query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return followedBy, nil
		}
		return nil, fmt.Errorf("PScaleSocials: unable to get followed by - %w", err)
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

func (sDB *SocialsDB) Follow(userID, toFollow int) error {
	tx, err := sDB.DB.Begin()
	if err != nil {
		return fmt.Errorf("PScaleSocials: unable to begin tx for following user - %w", err)
	}
	defer tx.Rollback()

	query := "INSERT INTO follows (user_id, follows) VALUES (?, ?)"

	_, err = tx.Exec(query, userID, toFollow)
	if err != nil {
		return fmt.Errorf("PScaleSocials: unable to add following relation to db - %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("PScaleSocials: unable to commit tx to add follow to db - %w", err)
	}

	return nil
}

func (sDB *SocialsDB) UnFollow(userID, toUnFollow int) error {
	tx, err := sDB.DB.Begin()
	if err != nil {
		return fmt.Errorf("PScaleSocials: unable to begin tx for unfollowing user - %w", err)
	}
	defer tx.Rollback()

	query := "DELETE FROM follows WHERE user_id = ? AND follows = ?"
	_, err = tx.Exec(query, userID, toUnFollow)
	if err != nil {
		return fmt.Errorf("PScaleSocials: unable to delete from follows table in db - %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("PScaleSocials: unable to commit tx to delete from follows table in db - %w", err)
	}

	return nil
}

func (sDB *SocialsDB) Like(userID, postID int) error {
	tx, err := sDB.DB.Begin()
	if err != nil {
		return fmt.Errorf("PScaleSocials: unable to begin tx for adding like - %w", err)
	}
	defer tx.Rollback()

	queryOnPosts := "UPDATE notes SET likes = likes + 1 WHERE id = ?"
	_, err = tx.Exec(queryOnPosts, postID)
	if err != nil {
		return fmt.Errorf("PScaleSocials: unable to add like to notes table in db - %w", err)
	}

	queryOnLikes := "INSERT INTO likes_list (post_id, liked_by) VALUES (?, ?)"
	_, err = tx.Exec(queryOnLikes, postID, userID)
	if err != nil {
		return fmt.Errorf("PScaleSocials: unable to add like to likes_list table in db - %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("PScaleSocials: unable to commit tx to add like to db - %w", err)
	}

	return nil
}

func (sDB *SocialsDB) Unlike(userID, postID int) error {
	tx, err := sDB.DB.Begin()
	if err != nil {
		return fmt.Errorf("PScaleSocials: unable to begin tx for adding like - %w", err)
	}
	defer tx.Rollback()

	queryOnPosts := "UPDATE notes SET likes = likes - 1 WHERE id = ?"
	_, err = tx.Exec(queryOnPosts, postID)
	if err != nil {
		return fmt.Errorf("PScaleSocials: unable to add like to notes table in db - %w", err)
	}

	queryOnLikes := "DELETE FROM likes_list WHERE post_id = ? AND liked_by = ?"
	_, err = tx.Exec(queryOnLikes, postID, userID)
	if err != nil {
		return fmt.Errorf("PScaleSocials: unable to add like to likes_list table in db - %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("PScaleSocials: unable to commit tx to add like to db - %w", err)
	}

	return nil
}
