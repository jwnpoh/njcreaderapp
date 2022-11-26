package broker

import (
	"net/http"
	"strconv"

	"github.com/jwnpoh/njcreaderapp/backend/services/serializer"
)

func (b *broker) GetFriends(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("userid")

	userID, err := strconv.Atoi(q)
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
