package broker

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
	"github.com/jwnpoh/njcreaderapp/backend/services/profanity"
	"github.com/jwnpoh/njcreaderapp/backend/services/serializer"
)

// GetArticle gets the article information to pass back to the frontend in order for the user to create a note.
func (b *broker) GetArticle(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("id")

	id, err := strconv.Atoi(q)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to parse article id requested", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	data, err := b.Articles.GetArticle(id)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to get article from database", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	err = data.Encode(w, http.StatusAccepted)
	if err != nil {
		data.ErrorJson(w, err)
		b.Logger.Error(data, r)
	}
}

// GetPost gets the post information to pass back to the frontend in order for the user to edit a note.
func (b *broker) GetPost(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("id")

	id, err := strconv.Atoi(q)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to parse article id requested", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	data, err := b.Posts.GetPost(id)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to get article from database", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	err = data.Encode(w, http.StatusAccepted)
	if err != nil {
		data.ErrorJson(w, err)
		b.Logger.Error(data, r)
	}
}

// GetPublicFeed gets all public notes of all users.
func (b *broker) GetPublicFeed(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("user")

	var data serializer.Serializer
	if q == "all" {
		d, err := b.Posts.GetAllPublicPosts()
		if err != nil {
			s := serializer.NewSerializer(true, "unable to get public posts from database", err)
			s.ErrorJson(w, err)
			b.Logger.Error(s, r)
			return
		}
		data = d
		err = data.Encode(w, http.StatusAccepted)
		if err != nil {
			data.ErrorJson(w, err)
			b.Logger.Error(data, r)
		}
		return
	}

	id, err := strconv.Atoi(q)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to parse user id requested", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	data, err = b.Posts.GetPublicPosts(id)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to get public posts from database", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	err = data.Encode(w, http.StatusAccepted)
	if err != nil {
		data.ErrorJson(w, err)
		b.Logger.Error(data, r)
	}
}

// GetFollowing gets all public notes of users that are followed by the user with the userID specified in the url query of the GET request.
func (b *broker) GetFollowing(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("user")

	id, err := strconv.Atoi(q)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to parse user id requested", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	data, err := b.Posts.GetFollowingPosts(id)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to get public posts from users followed", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	err = data.Encode(w, http.StatusAccepted)
	if err != nil {
		data.ErrorJson(w, err)
		b.Logger.Error(data, r)
	}
}

// GetNotebook gets all notes created by the user with the userID specified in the url query of the GET request.
func (b *broker) GetNotebook(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("user")

	id, err := strconv.Atoi(q)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to parse user id requested", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	data, err := b.Posts.GetOwnPosts(id)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to get public posts from users followed", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	err = data.Encode(w, http.StatusAccepted)
	if err != nil {
		data.ErrorJson(w, err)
		b.Logger.Error(data, r)
	}
}

// InsertPost inserts a post created in the frontend. Some processing is done first like profanity checking and date formatting before being sent to the posts service to be passed on to the DB.
func (b *broker) InsertPost(w http.ResponseWriter, r *http.Request) {
	var input *core.PostPayload

	err := serializer.NewSerializer(false, "", nil).Decode(w, r, &input)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to decode input data", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	if input.Public == "on" {
		toCheck := []string{input.TLDR, input.Examples, input.Notes, strings.Join(input.Tags, ",")}
		for _, v := range toCheck {
			profanityCheck := profanity.CheckProfanity(v)
			if profanityCheck.IsProfane {
				s := serializer.NewSerializer(true, fmt.Sprintf("Please use clean language when posting publicly on this platform.\nThe system auto-detected the use of the profanity: '%s'.\nIf this is a false positive, please report the false positive to the system admin via your teacher.", profanityCheck.Profanity), input.UserID)
				s.Encode(w, http.StatusBadRequest)
				b.Logger.Info(s, r)
				return
			}
		}
	}

	date := formatDate(input.Date)
	input.Date = date

	err = b.Posts.AddPost(input)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to add new post", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	err = serializer.NewSerializer(false, "successfully added post", nil).Encode(w, http.StatusAccepted)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to add post", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}
}

// DeletePost deletes posts specified in the slice of postIDs specified in the POST request.
func (b *broker) DeletePost(w http.ResponseWriter, r *http.Request) {
	var input int

	err := serializer.NewSerializer(false, "", nil).Decode(w, r, &input)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to decode input data", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	err = b.Posts.DeletePost(input)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to delete posts", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	err = serializer.NewSerializer(false, "successfully deleted posts", nil).Encode(w, http.StatusAccepted)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to delete posts", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}
}

func (b *broker) UpdatePost(w http.ResponseWriter, r *http.Request) {
	var input struct {
		PostID string            `json:"post_id"`
		Post   *core.PostPayload `json:"post"`
	}

	err := serializer.NewSerializer(false, "", nil).Decode(w, r, &input)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to decode input data", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	id, err := strconv.Atoi(input.PostID)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to convert post_id from string to int", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	if input.Post.Public == "on" {
		toCheck := []string{input.Post.TLDR, input.Post.Examples, input.Post.Notes, strings.Join(input.Post.Tags, ",")}
		for _, v := range toCheck {
			profanityCheck := profanity.CheckProfanity(v)
			if profanityCheck.IsProfane {
				s := serializer.NewSerializer(true, fmt.Sprintf("Please use clean language on this platform.\nThe system auto-detected the use of the profanity: '%s'.\nIf this is a false positive, please report the false positive via the feedback form.", profanityCheck.Profanity), input.Post.UserID)
				s.Encode(w, http.StatusBadRequest)
				b.Logger.Info(s, r)
				return
			}
		}
	}

	date := formatDate(input.Post.Date)
	input.Post.Date = date

	err = b.Posts.UpdatePost(id, input.Post)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to add new post", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	err = serializer.NewSerializer(false, "successfully added post", nil).Encode(w, http.StatusAccepted)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to add post", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}
}

// GetLikedPosts gets a slice of ints for the post ids of all the posts like by the user with the userID specified in the url query of the GET request.
func (b *broker) GetLikedPosts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("user")

	id, err := strconv.Atoi(q)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to parse user id requested", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	data, err := b.Posts.GetLikedPosts(id)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to get user's liked posts", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	err = data.Encode(w, http.StatusAccepted)
	if err != nil {
		data.ErrorJson(w, err)
		b.Logger.Error(data, r)
	}
}

// GetPostLikes gets the display names of users who have liked the post with postID specified in the url query of the GET request.
func (b *broker) GetPostLikes(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("post")

	id, err := strconv.Atoi(q)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to parse user id requested", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	data, err := b.Posts.GetLikes(id)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to get user's liked posts", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	err = data.Encode(w, http.StatusAccepted)
	if err != nil {
		data.ErrorJson(w, err)
		b.Logger.Error(data, r)
	}
}
