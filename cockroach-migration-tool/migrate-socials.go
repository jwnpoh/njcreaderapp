package main

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type follows struct {
	userName   string
	followName string
}

type likes struct {
	postAuthor       string
	postArticleTitle string
	likedBy          string
}

func migrateFollowsTable(pscaleDB, cockroachDB *sqlx.DB) error {
	fmt.Println("getting follows from pscale...")
	follows, err := getPscaleFollows(pscaleDB)
	if err != nil {
		return fmt.Errorf("PScaleFollows: unable to get follows from pscale - %w", err)
	}
	fmt.Printf("got %d follows from pscale\n", len(follows))

	fmt.Printf("attempting to insert %d follows from pscale to cockcroach...\n", len(follows))
	err = insertFollowsToCockroach(cockroachDB, follows)
	if err != nil {
		return fmt.Errorf("CockroachFollows: unable to insert follows to cockroach - %w\n", err)
	}

	return nil
}

func getPscaleFollows(pscaleDB *sqlx.DB) ([]follows, error) {
	series := make([]follows, 0)

	query := "SELECT user_id, follows FROM follows;"

	fmt.Println("running query")
	rows, err := pscaleDB.Queryx(query)
	if err != nil {
		return nil, fmt.Errorf("PScaleFollows: unable to query follows table - %w\n", err)
	}

	var count int
	for rows.Next() {
		var follow follows
		var userid, followID int
		err = rows.Scan(&userid, &followID)
		if err != nil {
			return nil, fmt.Errorf("PScaleFollows: error scanning row - %w\n", err)
		}
		fmt.Println("scanned follows table", userid, followID)

		query := "SELECT email FROM users WHERE id = ?"
		row := pscaleDB.QueryRow(query, userid)
		err := row.Scan(&follow.userName)
		if err != nil {
			return nil, fmt.Errorf("PScaleFollows: error scanning user email from pscale - %w\n", err)
		}
		fmt.Printf("the user email for id %d is %s\n", userid, follow.userName)

		row = pscaleDB.QueryRow(query, followID)
		err = row.Scan(&follow.followName)
		if err != nil {
			return nil, fmt.Errorf("PScaleFollows: error scanning user email from pscale - %w\n", err)
		}

		series = append(series, follow)
		count++
	}
	fmt.Printf("scanned a total of %d follows from pscale\n", len(series))

	return series, nil
}

func insertFollowsToCockroach(cockroachDB *sqlx.DB, follows []follows) error {
	if len(follows) < 1 {
		fmt.Println("did not receive follows to insert to cockroach.")
		return fmt.Errorf("did not receive follows to insert to cockroach.")
	}

	query := "INSERT INTO follows (user_id, follows) VALUES ($1, $2)"

	fmt.Println("running query")
	tx, err := cockroachDB.Begin()
	if err != nil {
		return fmt.Errorf("CockroachFollows: unable to begin tx for adding follows to cockroach - %w\n", err)
	}
	defer tx.Rollback()

	for i, user := range follows {
		var id string
		var userID, follow uuid.UUID
		// get uuids from cockroach
		err := cockroachDB.QueryRowx("SELECT id FROM users WHERE email = $1", user.userName).Scan(&id)
		if err != nil {
			return fmt.Errorf("CockroachFollows: unable to get uuids from cockroach\n")
		}
		userID, err = uuid.Parse(id)

		err = cockroachDB.QueryRowx("SELECT id FROM users WHERE email = $1", user.followName).Scan(&id)
		if err != nil {
			return fmt.Errorf("CockroachFollows: unable to get uuids from cockroach\n")
		}
		follow, err = uuid.Parse(id)

		fmt.Printf("found uuids for:\n%s: %s\n%s: %s\n", user.userName, userID, user.followName, follow)

		_, err = tx.Exec(query, userID, follow)
		if err != nil {
			return fmt.Errorf("CockroachFollows: unable to add follow for %s following %s to cockroach\n", user.followName, user.userName)
		}

		fmt.Printf("successfully inserted follow  for %s following %s\n", user.followName, user.userName)
		if i >= 0 && i%200 == 0 || i == len(follows)-1 {
			err = tx.Commit()
			if err != nil {
				return fmt.Errorf("CockroachFollows: unable to commit tx to db - %w", err)
			}
			tx, err = cockroachDB.Begin()
			if err != nil {
				return fmt.Errorf("CockroachFollows: unable to begin tx for adding follows input to cockroach - %w", err)
			}
			defer tx.Rollback()
		}
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("CockroachFollows: unable to commit tx to db - %w", err)
	}

	return nil
}

func migrateLikesTable(pscaleDB, cockroachDB *sqlx.DB) error {
	fmt.Println("getting likes from pscale...")
	likes, err := getPscaleLikes(pscaleDB)
	if err != nil {
		return fmt.Errorf("PScaleLikes: unable to get likes from pscale - %w", err)
	}
	fmt.Printf("got %d likes from pscale\n", len(likes))

	fmt.Printf("attempting to insert %d likes from pscale to cockcroach...\n", len(likes))
	err = insertLikesToCockroach(cockroachDB, likes)
	if err != nil {
		return fmt.Errorf("CockroachLikes: unable to insert likes to cockroach - %w\n", err)
	}

	return nil
}

func getPscaleLikes(pscaleDB *sqlx.DB) ([]likes, error) {
	series := make([]likes, 0)

	query := "SELECT post_id, liked_by FROM likes_list;"

	fmt.Println("running query")
	rows, err := pscaleDB.Queryx(query)
	if err != nil {
		return nil, fmt.Errorf("PScaleLikes: unable to query likes table - %w\n", err)
	}

	var count int
	for rows.Next() {
		var like likes
		var postID, likedBy int
		err = rows.Scan(&postID, &likedBy)
		if err != nil {
			return nil, fmt.Errorf("PScaleLikes: error scanning row - %w\n", err)
		}
		fmt.Println("scanned likes table", postID, likedBy)

		var postAuthor, postArticleTitle string
		query := "SELECT author, article_title FROM notes WHERE id = ?"
		row := pscaleDB.QueryRow(query, postID)
		err := row.Scan(&postAuthor, &postArticleTitle)
		if err != nil {
			return nil, fmt.Errorf("PScaleLikes: error scanning user email from pscale - %w\n", err)
		}
		like.postAuthor = postAuthor
		like.postArticleTitle = postArticleTitle
		fmt.Printf("the liked post id: %d is written by %s on %s\n", postID, like.postAuthor, like.postArticleTitle)

		query = "SELECT email FROM users WHERE id = ?"
		row = pscaleDB.QueryRow(query, likedBy)
		err = row.Scan(&like.likedBy)
		if err != nil {
			return nil, fmt.Errorf("PScaleLikes: error scanning user email from pscale - %w\n", err)
		}
		fmt.Printf("the liked post is liked by %s\n", like.likedBy)

		series = append(series, like)
		count++
	}

	fmt.Printf("scanned a total of %d likes from pscale\n", len(series))

	return series, nil
}

func insertLikesToCockroach(cockroachDB *sqlx.DB, likes []likes) error {
	if len(likes) < 1 {
		fmt.Println("did not receive likes to insert to cockroach.")
		return fmt.Errorf("did not receive likes to insert to cockroach.")
	}

	query := "INSERT INTO likes_list (post_id, liked_by) VALUES ($1, $2)"

	fmt.Println("running query")
	tx, err := cockroachDB.Begin()
	if err != nil {
		return fmt.Errorf("CockroachLikes: unable to begin tx for adding likes to cockroach - %w\n", err)
	}
	defer tx.Rollback()

	for i, like := range likes {
		var id string
		var postID, likedBy uuid.UUID
		// get uuids from cockroach
		err := cockroachDB.QueryRowx("SELECT id FROM notes WHERE author = $1 AND article_title = $2", like.postAuthor, like.postArticleTitle).Scan(&id)
		if err != nil {
			return fmt.Errorf("CockroachLikes: unable to get uuids from cockroach\n")
		}
		postID, err = uuid.Parse(id)

		err = cockroachDB.QueryRowx("SELECT id FROM users WHERE email = $1", like.likedBy).Scan(&id)
		if err != nil {
			return fmt.Errorf("CockroachLikes: unable to get uuids from cockroach\n")
		}
		likedBy, err = uuid.Parse(id)

		fmt.Printf("found uuids for the post %s liked by %s\n", postID, likedBy)

		_, err = tx.Exec(query, postID, likedBy)
		if err != nil {
			return fmt.Errorf("CockroachLikes: unable to add follow for %s following %s to cockroach\n", postID, likedBy)
		}

		fmt.Printf("successfully inserted follow  for %s following %s\n", postID, likedBy)
		if i >= 0 && i%200 == 0 || i == len(likes)-1 {
			err = tx.Commit()
			if err != nil {
				return fmt.Errorf("CockroachLikes: unable to commit tx to db - %w", err)
			}
			tx, err = cockroachDB.Begin()
			if err != nil {
				return fmt.Errorf("CockroachLikes: unable to begin tx for adding likes input to cockroach - %w", err)
			}
			defer tx.Rollback()
		}
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("CockroachLikes: unable to commit tx to db - %w", err)
	}

	return nil
}
