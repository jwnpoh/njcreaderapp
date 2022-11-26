package broker

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
	"github.com/jwnpoh/njcreaderapp/backend/services/serializer"
)

// Get makes a call to the articles service to retrieve articles from the db for a given page.
func (b *broker) Get(w http.ResponseWriter, r *http.Request) {
	page := chi.URLParam(r, "page")

	n, err := strconv.Atoi(page)
	if err != nil {
		s := serializer.NewSerializer(true, "something went wrong with the page number", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	data, err := b.Articles.Get(n, 8)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to get articles from database", err)
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

// Get100 makes a call to the articles service to retrieve 100 articles from the db for admin editing/deleting.
func (b *broker) Get100(w http.ResponseWriter, r *http.Request) {
	data, err := b.Articles.Get(1, 100)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to get articles from database", err)
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

// Find makes a call to the database via the articles service to search for a match of the given search term specified in the url params.
func (b *broker) Find(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("term")

	data, err := b.Articles.Find(q)
	if err != nil {
		s := serializer.NewSerializer(true, fmt.Sprintf("unable to find results for search term %s", q), err)
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

// Store parses the new article input in the request body and sends it to the db via articles service.
func (b *broker) Store(w http.ResponseWriter, r *http.Request) {
	input := make(core.ArticlePayload, 0)

	s := serializer.NewSerializer(false, "successfully stored articles", nil)
	err := s.Decode(w, r, &input)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to decode input data", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	for i, item := range input {
		date := formatDate(item.Date)
		input[i].Date = date
	}

	err = b.Articles.Store(input)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to store input data", err)
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

// Update parses the updated articles and sends the update to the db via the articles service.
func (b *broker) Update(w http.ResponseWriter, r *http.Request) {
	input := make(core.ArticlePayload, 0)

	s := serializer.NewSerializer(false, "successfully updated article", nil)
	err := s.Decode(w, r, &input)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to decode input data", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	for i, item := range input {
		date := formatDate(item.Date)
		input[i].Date = date
	}

	err = b.Articles.Update(input)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to store input data", err)
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

// Delete takes a range of article ids and deletes theme from the database.
func (b *broker) Delete(w http.ResponseWriter, r *http.Request) {
	input := make([]string, 0)

	s := serializer.NewSerializer(false, "successfully deleted articles", nil)
	err := s.Decode(w, r, &input)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to decode input data", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	err = b.Articles.Delete(input)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to store input data", err)
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
