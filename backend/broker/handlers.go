package broker

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
	"github.com/jwnpoh/njcreaderapp/backend/serializer"
	"github.com/jwnpoh/njcreaderapp/backend/services/articles"
)

// Get makes a call to the articles service to retrieve articles from the db for a given page.
func (b *broker) Get(w http.ResponseWriter, r *http.Request) {
	page := chi.URLParam(r, "page")

	n, err := strconv.Atoi(page)
	if err != nil {
		s := serializer.NewSerializer(true, "something went wrong with the page number", nil)
		s.ErrorJson(w, err)
		b.Logger.Error(r.Method, r.URL, s)
	}

	service, err := articles.NewArticlesService()
	if err != nil {
		s := serializer.NewSerializer(true, "unable to start articles service", nil)
		s.ErrorJson(w, err)
		b.Logger.Error(r.Method, r.URL, s)
	}

	data, err := service.Get(n)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to get articles from database", nil)
		s.ErrorJson(w, err)
		b.Logger.Error(r.Method, r.URL, err)
	}

	b.Logger.Error(data.Encode(w, http.StatusAccepted))
	b.Logger.Success(r.Method, r.URL)
}

// Store parses the new article input in the request body and sends it to the db via articles service.
func (b *broker) Store(w http.ResponseWriter, r *http.Request) {
	data := make(core.ArticleSeries, 0)

	s := serializer.NewSerializer(false, "", nil)
	err := s.Decode(w, r, &data)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to decode input data", nil)
		s.ErrorJson(w, err)
		b.Logger.Error(r.Method, r.URL, s)
	}

	service, err := articles.NewArticlesService()
	if err != nil {
		s := serializer.NewSerializer(true, "unable to start articles service to store input", nil)
		s.ErrorJson(w, err)
		b.Logger.Error(r.Method, r.URL, s)

	}

	err = service.Store(data)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to store input data", nil)
		s.ErrorJson(w, err)
		b.Logger.Error(r.Method, r.URL, s)
	}
	b.Logger.Success(r.Method, r.URL)
}

// Find makes a call to the database via the articles service to search for a match of the given search term specified in the url params.
func (b *broker) Find(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("term")
	service, err := articles.NewArticlesService()
	if err != nil {
		s := serializer.NewSerializer(true, "unable to start articles service for search", nil)
		s.ErrorJson(w, err)
		b.Logger.Error(r.Method, r.URL, s)
	}

	data, err := service.Find(q)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to start find results for search term", nil)
		s.ErrorJson(w, err)
		b.Logger.Error(r.Method, r.URL, s)
	}

	data.Encode(w, http.StatusAccepted)
	b.Logger.Success(r.Method, r.URL)
}
