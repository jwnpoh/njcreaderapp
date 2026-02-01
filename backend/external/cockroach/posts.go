package cockroach

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
)

type PostsDB struct {
	DB *sqlx.DB
}

// NewCockroachDB returns a connection interface for the application to connect to the planetscale database.
func NewPostsDB(dsn string) (*PostsDB, error) {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("CockroachPosts: unable to initialize Cockroach database - %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("CockroachPosts: no response from cockroach database - %w", err)
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)

	return &PostsDB{DB: db}, nil
}

func (pDB *PostsDB) GetAllPublicPosts() (*core.Posts, error) {
	posts := make(core.Posts, 0)

	query := "SELECT * FROM notes WHERE public = true ORDER BY created_at DESC"

	rows, err := pDB.DB.Queryx(query)
	if err != nil {
		return nil, fmt.Errorf("CockroachPost-GetAllPublicPosts: unable to query notes table - %w", err)
	}

	for rows.Next() {
		var post core.Post
		var tags, author, authorClass sql.NullString
		err = rows.Scan(&post.ID, &post.UserID, &author, &authorClass, &post.Likes, &post.TLDR, &post.Examples, &post.Notes, &tags, &post.CreatedAt, &post.Public, &post.ArticleID, &post.ArticleTitle, &post.ArticleURL)
		if err != nil {
			return nil, fmt.Errorf("CockroachPosts-GetAllPublicPosts(): error scanning row - %w", err)
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

		posts = append(posts, post)
	}

	return &posts, nil
}

func (pDB *PostsDB) GetPosts(userIDs []uuid.UUID, public bool) (*core.Posts, error) {
	posts := make(core.Posts, 0)

	query := parseQuery(userIDs)

	if public {
		query += " AND public = true"
	}

	query += " ORDER BY created_at DESC"

	rows, err := pDB.DB.Queryx(query)
	if err != nil {
		return nil, fmt.Errorf("CockroachPosts-GetPosts: unable to query notes table with query %s - %w", query, err)
	}

	for rows.Next() {
		var post core.Post
		var tags, author, authorClass sql.NullString
		err = rows.Scan(&post.ID, &post.UserID, &author, &authorClass, &post.Likes, &post.TLDR, &post.Examples, &post.Notes, &tags, &post.CreatedAt, &post.Public, &post.ArticleID, &post.ArticleTitle, &post.ArticleURL)
		if err != nil {
			return nil, fmt.Errorf("CockroachPosts-GetPosts: error scanning row - %w", err)
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

		posts = append(posts, post)
	}

	return &posts, nil
}

func (pDB *PostsDB) AddPost(post *core.Post) error {
	query := "INSERT INTO notes (user_id, author, author_class, likes, tldr, examples, notes, tags, created_at, public, article_id, article_title, article_url) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) ON CONFLICT DO NOTHING"

	tx, err := pDB.DB.Begin()
	if err != nil {
		return fmt.Errorf("CockroachPosts-AddPost: unable to begin tx for adding post to db - %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(query, post.UserID, post.Author, post.AuthorClass, post.Likes, post.TLDR, post.Examples, post.Notes, strings.Join(post.Tags, ","), post.CreatedAt, post.Public, post.ArticleID, post.ArticleTitle, post.ArticleURL)
	if err != nil {
		return fmt.Errorf("CockroachPosts-AddPost: unable to add post to db - %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("CockroachPosts-AddPost: unable to commit tx to add post to db - %w", err)
	}

	return nil
}

func (pDB *PostsDB) GetPost(id uuid.UUID) (*core.Post, error) {
	var post core.Post

	query := "SELECT * FROM notes WHERE id = $1"

	row := pDB.DB.QueryRowx(query, id)

	var tags, author, authorClass sql.NullString
	err := row.Scan(&post.ID, &post.UserID, &author, &authorClass, &post.Likes, &post.TLDR, &post.Examples, &post.Notes, &tags, &post.CreatedAt, &post.Public, &post.ArticleID, &post.ArticleTitle, &post.ArticleURL)
	if err != nil {
		return nil, fmt.Errorf("CockroachPosts-GetPost: error scanning row - %w", err)
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

	return &post, nil
}

func (pDB *PostsDB) DeletePost(postID uuid.UUID) error {
	query := fmt.Sprintf("DELETE FROM notes WHERE id = %d", postID)

	tx, err := pDB.DB.Begin()
	if err != nil {
		return fmt.Errorf("CockroachPosts-DeletePost: unable to begin tx for deleting notes from db - %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(query)
	if err != nil {
		return fmt.Errorf("CockroachPosts-DeletePost: unable to delete notes from db - %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("CockroachCockroachosts: unable to commit tx to delete notes from db - %w", err)
	}

	return nil
}

func (pDB *PostsDB) UpdatePost(postID uuid.UUID, post *core.Post) error {
	query := "UPDATE notes SET user_id = $1, author = $2, author_class = $3, likes = $4, tldr = $5, examples = $6, notes = $7, tags = $8, created_at = $9, public = $10, article_id = $11, article_title = $12, article_url = $13 WHERE id = $14"

	tx, err := pDB.DB.Begin()
	if err != nil {
		return fmt.Errorf("CockroachPosts-UpdatePost: unable to begin tx for adding post to db - %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(query, post.UserID, post.Author, post.AuthorClass, post.Likes, post.TLDR, post.Examples, post.Notes, strings.Join(post.Tags, ","), post.CreatedAt, post.Public, post.ArticleID, post.ArticleTitle, post.ArticleURL, postID)
	if err != nil {
		return fmt.Errorf("CockroachPosts-UpdatePost: unable to add post to db - %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("CockroachPosts-UpdatePost: unable to commit tx to add post to db - %w", err)
	}

	return nil
}

func (pDB *PostsDB) GetLikes(id uuid.UUID, userOrPost string) ([]uuid.UUID, error) {
	likedBys := make([]uuid.UUID, 0)

	var query string

	switch userOrPost {
	case "post":
		query = "SELECT liked_by FROM likes_list WHERE post_id = $1"
	case "user":
		query = "SELECT post_id FROM likes_list WHERE liked_by = $1"
	default:
		return nil, fmt.Errorf("CockroachPosts-GetLikes: error interpreting 'user' or 'post'")
	}

	rows, err := pDB.DB.Queryx(query, id)
	if err != nil {
		return nil, fmt.Errorf("CockroachPosts-GetLikes: unable to query likes_list table - %w", err)
	}

	var resID uuid.UUID
	for rows.Next() {
		err = rows.Scan(&resID)
		if err != nil {
			return nil, fmt.Errorf("CockroachPosts-GetLikes: error scanning row - %w", err)
		}
		likedBys = append(likedBys, resID)
	}

	return likedBys, nil
}

func parseQuery(userIDs []uuid.UUID) string {
	if len(userIDs) == 1 {
		return fmt.Sprintf("SELECT * FROM notes WHERE user_id = '%s'", userIDs[0])
	}

	query := strings.Builder{}
	query.WriteString("SELECT * FROM notes WHERE ")
	for i, v := range userIDs {
		if i < len(userIDs)-1 {
			query.WriteString(fmt.Sprintf("user_id = '%s' OR ", v))
			continue
		}
		query.WriteString(fmt.Sprintf("user_id = '%s'", v))
	}

	return query.String()
}
