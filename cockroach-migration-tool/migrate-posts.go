package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jwnpoh/njcreaderapp/cockroach-migration-tool/core"
)

func migrateNotesTable(pscaleDB, cockroachDB *sqlx.DB) error {
	fmt.Println("getting notes from pscale...")
	notes, err := getPscaleNotes(pscaleDB)
	if err != nil {
		return fmt.Errorf("PScaleNotes: unable to get Notes from pscale - %w", err)
	}
	fmt.Printf("got %d Notes from pscale\n", len(notes))

	fmt.Printf("attempting to insert %d notes from pscale to cockcroach...\n", len(notes))
	err = insertNotesToCockroach(cockroachDB, notes)
	if err != nil {
		return fmt.Errorf("CockroachNotes: unable to insert notes to cockroach - %w\n", err)
	}

	return nil
}

func getPscaleNotes(pscaleDB *sqlx.DB) (core.Posts, error) {
	series := make(core.Posts, 0)

	query := "SELECT * FROM notes"

	fmt.Println("running query")
	rows, err := pscaleDB.Queryx(query)
	if err != nil {
		return nil, fmt.Errorf("PScaleNotes: unable to query notes table - %w\n", err)
	}

	var count int
	for rows.Next() {
		var post core.Post
		var id, oldUserID, oldArticleID int
		var tags, author, authorClass sql.NullString
		err = rows.Scan(&id, &oldUserID, &author, &authorClass, &post.Likes, &post.TLDR, &post.Examples, &post.Notes, &tags, &post.CreatedAt, &post.Public, &oldArticleID, &post.ArticleTitle, &post.ArticleURL)
		if err != nil {
			return nil, fmt.Errorf("PScaleNotes: error scanning row - %w\n", err)
		}
		if tags.Valid {
			post.Tags = strings.Split(tags.String, ",")
		}
		if author.Valid {
			post.Author = author.String
		}
		if authorClass.Valid {
			post.AuthorClass = authorClass.String
		}
		post.Date = time.Unix(post.CreatedAt, 0).Format("Jan 2, 2006 15:04:05")

		series = append(series, post)
		count++
		fmt.Printf("scanned %d notes...\n", count)
	}

	fmt.Printf("scanned a total of %d notes from pscale\n", len(series))

	return series, nil
}

func insertNotesToCockroach(cockroachDB *sqlx.DB, notes core.Posts) error {
	if len(notes) < 1 {
		fmt.Println("did not receive notes to insert to cockroach.")
		return fmt.Errorf("did not receive notes to insert to cockroach.")
	}

	query := "INSERT INTO notes (user_id, author, author_class, likes, tldr, examples, notes, tags, created_at, public, article_id, article_title, article_url) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) ON CONFLICT DO NOTHING"

	fmt.Println("running query")
	tx, err := cockroachDB.Begin()
	if err != nil {
		return fmt.Errorf("Cockroach: unable to begin tx for adding notes input to cockroach - %w\n", err)
	}
	defer tx.Rollback()

	for i, note := range notes {
		tx, err := cockroachDB.Begin()
		if err != nil {
			return fmt.Errorf("CockroachPosts: unable to begin tx for adding post to db - %w", err)
		}
		defer tx.Rollback()

		// get user uuid from cockroach as id is now uuid and not int
		newUserID, err := getNewUserID(cockroachDB, note.Author, note.AuthorClass)
		if err != nil {
			return fmt.Errorf("CockroachNotes: can't get newUserID - %w\n", err)
		}
		note.UserID = newUserID

		// get article uuid as id is now uuid and not int
		newArticleID, err := getNewArticleID(cockroachDB, note.ArticleURL)
		if err != nil {
			return fmt.Errorf("CockroachNotes: can't get newArticleID - %w\n", err)
		}
		note.ArticleID = newArticleID

		// insert note to cockroach
		_, err = tx.Exec(query, note.UserID, note.Author, note.AuthorClass, note.Likes, note.TLDR, note.Examples, note.Notes, strings.Join(note.Tags, ","), note.CreatedAt, note.Public, note.ArticleID, note.ArticleTitle, note.ArticleURL)
		if err != nil {
			return fmt.Errorf("CockroachPosts: unable to add post to db - %w", err)
		}

		fmt.Printf("successfully inserted note #%d \n", i+1)
		if i == len(notes)-1 {
			err = tx.Commit()
			if err != nil {
				return fmt.Errorf("CockroachNotes: unable to commit tx to insert notes to db - %w", err)
			}
		}
	}

	return nil
}

func getNewUserID(cockroachDB *sqlx.DB, author, authorClass string) (uuid.UUID, error) {
	fmt.Printf("getting new user uuid from cockroach for %s from %s\n", author, authorClass)

	query := "SELECT id FROM users WHERE display_name = $1 AND class = $2"
	var uuid uuid.UUID

	row := cockroachDB.QueryRowx(query, author, authorClass)
	err := row.Scan(&uuid)
	if err != nil {
		return uuid, fmt.Errorf("CockroachUsers: error scanning row - %w", err)
	}

	fmt.Printf("%s's uuid is %s", author, uuid)
	return uuid, nil
}

func getNewArticleID(cockroachDB *sqlx.DB, articleURL string) (uuid.UUID, error) {
	fmt.Printf("getting new article uuid from cockroach for %s\n", articleURL)

	query := "SELECT id FROM articles WHERE url = $1"
	var uuid uuid.UUID

	row := cockroachDB.QueryRowx(query, articleURL)
	err := row.Scan(&uuid)
	if err != nil {
		return uuid, fmt.Errorf("CockroachUsers: error scanning row - %w", err)
	}

	fmt.Printf("The article uuid is %s", uuid)
	return uuid, nil
}
