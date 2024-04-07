package broker

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jwnpoh/njcreaderapp/backend/services/serializer"
)

func (b *broker) GetFriends(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("user")

	userID, err := uuid.Parse(q)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to get parse user id from request", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	data, err := b.Socials.GetFriends(userID)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to get friends", err)
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

func (b *broker) Follow(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID   uuid.UUID `json:"user_id"`
		ToFollow uuid.UUID `json:"to_follow"`
		Follow   bool      `json:"follow"`
	}

	s := serializer.NewSerializer(false, "", nil)
	s.Decode(w, r, &input)

	data, err := b.Socials.Follow(input.UserID, input.ToFollow, input.Follow)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to follow user", err)
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

func (b *broker) Like(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID uuid.UUID `json:"user_id"`
		PostID uuid.UUID `json:"post_id"`
		Like   bool      `json:"like"`
	}

	s := serializer.NewSerializer(false, "", nil)
	s.Decode(w, r, &input)

	data, err := b.Socials.Like(input.UserID, input.PostID, input.Like)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to update like", err)
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
