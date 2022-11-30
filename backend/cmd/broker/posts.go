package broker

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
	"github.com/jwnpoh/njcreaderapp/backend/services/profanity"
	"github.com/jwnpoh/njcreaderapp/backend/services/serializer"
)

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

func (b *broker) InsertPost(w http.ResponseWriter, r *http.Request) {
	var input *core.PostPayload

	err := serializer.NewSerializer(false, "", nil).Decode(w, r, &input)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to decode input data", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	toCheck := []string{input.TLDR, input.Examples, input.Notes}
	for _, v := range toCheck {
		profanityCheck := profanity.CheckProfanity(v)
		if profanityCheck.IsProfane {
			s := serializer.NewSerializer(true, fmt.Sprintf("Please use clean language on this platform.\nThe system auto-detected the use of the profanity: '%s'.\nIf this is a false positive, please report the false positive via the feedback form.", profanityCheck.Profanity), input.UserID)
			s.Encode(w, http.StatusBadRequest)
			b.Logger.Info(s, r)
			return
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

func (b *broker) DeletePost(w http.ResponseWriter, r *http.Request) {
	var input []string

	s := serializer.NewSerializer(false, "", nil)
	err := s.Decode(w, r, &input)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to decode input data", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	err = b.Posts.DeletePosts(input)
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
