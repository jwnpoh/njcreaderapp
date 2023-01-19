package broker

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
	"github.com/jwnpoh/njcreaderapp/backend/services/serializer"
)

func (b *broker) GetLongTopics(w http.ResponseWriter, r *http.Request) {
	data, err := b.Longs.GetTopics()
	if err != nil {
		s := serializer.NewSerializer(true, "unable to get topics for long articles", err)
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

func (b *broker) GetLong(w http.ResponseWriter, r *http.Request) {
	topic := chi.URLParam(r, "topic")

	data, err := b.Longs.Get(topic)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to get topics for long articles", err)
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

func (b *broker) StoreLong(w http.ResponseWriter, r *http.Request) {
	var input string

	s := serializer.NewSerializer(false, "successfully stored long articles", nil)
	err := s.Decode(w, r, &input)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to decode input data", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	fmt.Println(input)

	s, err = b.Longs.Store(input)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to store input data in database", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	err = s.Encode(w, http.StatusAccepted)
	if err != nil {
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
	}
}

func (b *broker) UpdateLong(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		var input core.Long

		s := serializer.NewSerializer(false, "successfully updated long article", nil)
		err := s.Decode(w, r, &input)
		if err != nil {
			s := serializer.NewSerializer(true, "unable to decode input data", err)
			s.ErrorJson(w, err)
			b.Logger.Error(s, r)
			return
		}

		s, err = b.Longs.Update(&input)
		if err != nil {
			s := serializer.NewSerializer(true, "unable to update input data in database", err)
			s.ErrorJson(w, err)
			b.Logger.Error(s, r)
			return
		}

		err = s.Encode(w, http.StatusAccepted)
		if err != nil {
			s.ErrorJson(w, err)
			b.Logger.Error(s, r)
		}
	} else if r.Method == "GET" {
		topic := "all"
		data, err := b.Longs.Get(topic)
		if err != nil {
			s := serializer.NewSerializer(true, "unable to get topics for long articles", err)
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
}

func (b *broker) DeleteLong(w http.ResponseWriter, r *http.Request) {

	input := make([]string, 0)

	s := serializer.NewSerializer(false, "successfully deleted long articles", nil)
	err := s.Decode(w, r, &input)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to decode input data", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	s, err = b.Longs.Delete(input)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to delete long articles", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	err = s.Encode(w, http.StatusAccepted)
	if err != nil {
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
	}
}