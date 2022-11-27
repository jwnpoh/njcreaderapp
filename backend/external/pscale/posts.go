package pscale

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
)

type PostsDB struct {
	DB *sqlx.DB
}

// NewPscaleDB returns a connection interface for the application to connect to the planetscale database.
func NewPostsDB(dsn string) (*PostsDB, error) {
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

	return &PostsDB{DB: db}, nil
}

func (pDB *PostsDB) GetAllPublicPosts() (*core.Posts, error) {
	posts := make(core.Posts, 0)

	query := "SELECT * FROM posts WHERE public = true ORDER BY created_at DESC"

	rows, err := pDB.DB.Queryx(query)
	if err != nil {
		return nil, fmt.Errorf("PScalePosts: unable to query posts table - %w", err)
	}

	for rows.Next() {
		var post core.Post
		var tags, author sql.NullString
		err = rows.Scan(&post.ID, &post.UserID, &author, &post.Likes, &post.TLDR, &post.Examples, &post.Notes, &tags, &post.CreatedAt, &post.Public, &post.ArticleID, &post.ArticleTitle, &post.ArticleURL)
		if err != nil {
			return nil, fmt.Errorf("PScalePosts: error scanning row - %w", err)
		}
		if tags.Valid {
			post.Tags = strings.Split(tags.String, ",")
		}
		if author.Valid {
			post.Author = author.String
		}
		post.Date = time.Unix(post.CreatedAt, 0).Format("Jan 2, 2006 15:04:05")

		posts = append(posts, post)
	}

	return &posts, nil
}

func (pDB *PostsDB) GetPosts(userIDs []int, public bool) (*core.Posts, error) {
	posts := make(core.Posts, 0)

	query := parseQuery(userIDs)

	if public {
		query += " AND public = true"
	}

	rows, err := pDB.DB.Queryx(query)
	if err != nil {
		return nil, fmt.Errorf("PScalePosts: unable to query posts table - %w", err)
	}

	for rows.Next() {
		var post core.Post
		var tags, author sql.NullString
		err = rows.Scan(&post.ID, &post.UserID, &author, &post.Likes, &post.TLDR, &post.Examples, &post.Notes, &tags, &post.CreatedAt, &post.Public, &post.ArticleID, &post.ArticleTitle, &post.ArticleURL)
		if err != nil {
			return nil, fmt.Errorf("PScalePosts: error scanning row - %w", err)
		}
		if tags.Valid {
			post.Tags = strings.Split(tags.String, ",")
		}
		if author.Valid {
			post.Author = author.String
		}
		post.Date = time.Unix(post.CreatedAt, 0).Format("Jan 2, 2006 15:04:05")

		posts = append(posts, post)
	}

	return &posts, nil
}

func (pDB *PostsDB) AddPost(post *core.Post) error {

	query := "INSERT INTO posts (user_id, author, likes, tldr, examples, notes, tags, created_at, public, article_id, article_title, article_url) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	tx, err := pDB.DB.Begin()
	if err != nil {
		return fmt.Errorf("PScalePosts: unable to begin tx for adding post to db - %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(query, post.UserID, post.Author, post.Likes, post.TLDR, post.Examples, post.Notes, strings.Join(post.Tags, ","), post.CreatedAt, post.Public, post.ArticleID, post.ArticleTitle, post.ArticleURL)
	if err != nil {
		return fmt.Errorf("PScalePosts: unable to add post to db - %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("PScalePosts: unable to commit tx to add post to db - %w", err)
	}

	return nil
}

func (pDB *PostsDB) DeletePosts(postIDs string) error {
	query := fmt.Sprintf("DELETE FROM posts WHERE id in (%s)", postIDs)

	tx, err := pDB.DB.Begin()
	if err != nil {
		return fmt.Errorf("PScalePosts: unable to begin tx for deleting posts from db - %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(query)
	if err != nil {
		return fmt.Errorf("PScalePosts: unable to delete posts from db - %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("PScalePosts: unable to commit tx to delete posts from db - %w", err)
	}

	return nil
}

func (pDB *PostsDB) GetLikes(postID int) ([]int, error) {
	likedBys := make([]int, 0)

	query := "SELECT liked_by FROM likes_list WHERE post_id = ?"

	rows, err := pDB.DB.Queryx(query, postID)
	if err != nil {
		return nil, fmt.Errorf("PScalePosts: unable to query likes_list table - %w", err)
	}

	var id int
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("PScalePosts: error scanning row - %w", err)
		}
		likedBys = append(likedBys, id)
	}

	return likedBys, nil
}

func parseQuery(userIDs []int) string {
	if len(userIDs) == 1 {
		return fmt.Sprintf("SELECT * FROM posts WHERE user_id = %d ORDER BY created_at DESC", userIDs[0])
	}

	params := strings.Builder{}
	for i, v := range userIDs {
		if i < len(userIDs)-1 {
			params.WriteString(fmt.Sprintf("%d, ", v))
			continue
		}
		params.WriteString(fmt.Sprintf("%d", v))
	}

	query := strings.Builder{}
	query.WriteString(fmt.Sprintf("SELECT * FROM posts WHERE user_id IN (%s) ORDER BY created_at DESC", params))

	return query.String()
}
