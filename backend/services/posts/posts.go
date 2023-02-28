package posts

import (
	"fmt"
	"strings"
	"time"

	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
	"github.com/jwnpoh/njcreaderapp/backend/services/articles"
	"github.com/jwnpoh/njcreaderapp/backend/services/serializer"
	"github.com/jwnpoh/njcreaderapp/backend/services/socials"
	"github.com/jwnpoh/njcreaderapp/backend/services/users"
)

type PostsDB interface {
	GetAllPublicPosts() (*core.Posts, error)
	GetPosts(userIDs []int, public bool) (*core.Posts, error)
	AddPost(post *core.Post) error
	GetPost(postID int) (*core.Post, error)
	DeletePost(postID int) error
	UpdatePost(postID int, post *core.Post) error
	GetLikes(id int, userOrPost string) ([]int, error)
}

type Posts struct {
	db PostsDB
	articles.ArticlesDB
	socials.SocialsDB
	users.UsersDB
}

func NewPostsDB(postsDB PostsDB, articlesDB articles.ArticlesDB, socialsDB socials.SocialsDB, usersDB users.UsersDB) *Posts {
	return &Posts{
		db:         postsDB,
		ArticlesDB: articlesDB,
		SocialsDB:  socialsDB,
		UsersDB:    usersDB,
	}
}

func (pDB *Posts) GetAllPublicPosts() (serializer.Serializer, error) {
	posts, err := pDB.db.GetAllPublicPosts()
	if err != nil {
		return nil, fmt.Errorf("unable to get posts from db - %w", err)
	}

	if len(*posts) < 1 {
		return serializer.NewSerializer(true, "There are currently no public notes from all users. Check back again later, or create your own note now.", nil), nil
	}

	var data core.Posts
	for _, v := range *posts {
		user, err := pDB.GetUser("id", v.UserID)
		if err != nil {
			return nil, fmt.Errorf("unable to get author info for notes from db - %w", err)
		}
		v.Author = user.DisplayName

		data = append(data, v)
	}

	return serializer.NewSerializer(false, "got all public notes", data), nil
}

func (pDB *Posts) GetPublicPosts(userID int) (serializer.Serializer, error) {
	params := []int{userID}

	posts, err := pDB.db.GetPosts(params, true)
	if err != nil {
		return nil, fmt.Errorf("unable to get notes from db - %w", err)
	}

	if len(*posts) < 1 {
		return serializer.NewSerializer(true, "This user currently has no public notes.", nil), nil
	}

	var data core.Posts
	for _, v := range *posts {
		user, err := pDB.GetUser("id", v.UserID)
		if err != nil {
			return nil, fmt.Errorf("unable to get author info for notes from db - %w", err)
		}
		v.Author = user.DisplayName

		data = append(data, v)
	}

	return serializer.NewSerializer(false, "got all user notes", data), nil
}

func (pDB *Posts) GetFollowingPosts(userID int) (serializer.Serializer, error) {
	params, err := pDB.GetFollowing(userID)
	if err != nil {
		return nil, fmt.Errorf("unable to get following - %w", err)
	}

	if len(params) < 1 {
		return serializer.NewSerializer(true, "You are currently not following anyone. Head over to the Discover page to find some people to follow!", nil), nil
	}

	posts, err := pDB.db.GetPosts(params, true)
	if err != nil {
		return nil, fmt.Errorf("unable to get notes from db - %w", err)
	}

	var data core.Posts
	for _, v := range *posts {
		user, err := pDB.GetUser("id", v.UserID)
		if err != nil {
			return nil, fmt.Errorf("unable to get author info for notes from db - %w", err)
		}
		v.Author = user.DisplayName

		data = append(data, v)
	}

	return serializer.NewSerializer(false, "There are currently no notes from the people that you are following. Try again later, or create your own note now.", data), nil
}

func (pDB *Posts) GetOwnPosts(userID int) (serializer.Serializer, error) {
	params := []int{userID}

	posts, err := pDB.db.GetPosts(params, false)
	if err != nil {
		return nil, fmt.Errorf("unable to get notes from db - %w", err)
	}

	if len(*posts) < 1 {
		return serializer.NewSerializer(true, "You have not created any note yet.", nil), nil
	}

	var data core.Posts
	for _, v := range *posts {
		user, err := pDB.GetUser("id", v.UserID)
		if err != nil {
			return nil, fmt.Errorf("unable to get author info for notes from db - %w", err)
		}
		v.Author = user.DisplayName

		data = append(data, v)
	}

	return serializer.NewSerializer(false, "got all user notes", data), nil
}

func (pDB *Posts) AddPost(post *core.PostPayload) error {
	newPost, err := parseNewPost(post)
	if err != nil {
		return fmt.Errorf("unable to parse input for new note - %w", err)
	}

	author, err := pDB.GetUser("id", newPost.UserID)
	if err != nil {
		return fmt.Errorf("unable to get author for new note - %w", err)
	}
	newPost.Author = author.DisplayName
	newPost.AuthorClass = author.Class

	err = pDB.db.AddPost(newPost)
	if err != nil {
		return fmt.Errorf("unable to add note - %w", err)
	}

	return nil
}

func (pDB *Posts) GetPost(postID int) (serializer.Serializer, error) {
	// get from Planetscale
	post, err := pDB.db.GetPost(postID)
	if err != nil {
		return nil, fmt.Errorf("unable to get note with id %d - %w", postID, err)
	}

	return serializer.NewSerializer(false, "got note", post), nil
}

func (pDB *Posts) DeletePost(postID int) error {
	// send to Planetscale
	err := pDB.db.DeletePost(postID)
	if err != nil {
		return fmt.Errorf("unable to delete note with id %d - %w", postID, err)
	}
	return nil
}

func (pDB *Posts) UpdatePost(postID int, post *core.PostPayload) error {
	newPost, err := parseNewPost(post)
	if err != nil {
		return fmt.Errorf("unable to parse input for new note - %w", err)
	}

	author, err := pDB.GetUser("id", newPost.UserID)
	if err != nil {
		return fmt.Errorf("unable to get author for new note - %w", err)
	}
	newPost.Author = author.DisplayName
	newPost.AuthorClass = author.Class

	err = pDB.db.UpdatePost(postID, newPost)
	if err != nil {
		return fmt.Errorf("unable to add note - %w", err)
	}

	return nil
}

func (pDB *Posts) GetLikes(postID int) (serializer.Serializer, error) {
	likedByIDs, err := pDB.db.GetLikes(postID, "post")
	if err != nil {
		return serializer.NewSerializer(true, fmt.Sprintf("unable to retrieve likes for note id %d - %v", postID, err), nil), err
	}

	likedByUsers := make([]string, 0, len(likedByIDs))
	for _, v := range likedByIDs {
		user, err := pDB.GetUser("id", v)
		if err != nil {
			return serializer.NewSerializer(true, fmt.Sprintf("unable to retrieve users who like note id %d - %v", postID, err), nil), err
		}
		likedByUsers = append(likedByUsers, user.DisplayName)
	}

	return serializer.NewSerializer(false, "successfully retrieved post likes", likedByUsers), nil
}

func (pDB *Posts) GetLikedPosts(userID int) (serializer.Serializer, error) {
	likedPosts, err := pDB.db.GetLikes(userID, "user")
	if err != nil {
		return serializer.NewSerializer(true, fmt.Sprintf("unable to retrieve likes by user id %d - %v", userID, err), nil), err
	}

	return serializer.NewSerializer(false, "successfully retrieved posts liked by user", likedPosts), nil
}

func parseNewPost(post *core.PostPayload) (*core.Post, error) {
	newPost := &core.Post{
		TLDR:         post.TLDR,
		Examples:     post.Examples,
		Notes:        post.Notes,
		UserID:       post.UserID,
		Likes:        post.Likes,
		ArticleID:    post.ArticleID,
		ArticleTitle: post.ArticleTitle,
		ArticleURL:   post.ArticleURL,
	}

	tags := parsePostTags(post.Tags)
	newPost.Tags = tags

	var public bool
	if post.Public == "on" {
		public = true
	}
	newPost.Public = public

	date := time.Now().Unix()
	newPost.CreatedAt = date

	return newPost, nil
}

func parsePostTags(input []string) []string {
	xinput := make([]string, 0, len(input))
	for _, v := range input {
		preTags := strings.Split(v, ",")
		xinput = append(xinput, preTags...)
	}

	tags := make([]string, 0, len(input))

	listOfTags := make(map[string]bool)

	for _, v := range xinput {
		tag := strings.TrimSpace(v)
		tag = strings.ReplaceAll(tag, "#", "")
		tag = strings.ToLower(tag)

		if !listOfTags[tag] {
			listOfTags[tag] = true
			tags = append(tags, tag)
		}
	}

	return tags
}