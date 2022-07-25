package broker

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jwnpoh/njcreaderapp/backend/serializer"
	"github.com/jwnpoh/njcreaderapp/backend/services/articles"
)

// GetLatest will be the "index" page of the application api.
func (b *broker) GetLatest(w http.ResponseWriter, r *http.Request) {
	serializer, err := articles.NewArticlesService().Get()
	if err != nil {
		serializer.ErrorJson(w, err)
	}

	serializer.Encode(w, http.StatusAccepted)
}

// GetPage will serve paged data if URL params are specified.
func (b *broker) GetPage(w http.ResponseWriter, r *http.Request) {
	pager := chi.URLParam(r, "pager")
	page, err := strconv.Atoi(pager)
	if err != nil {
		serializer.NewSerializer(true, "unable to parse page", nil).ErrorJson(w, err)
	}

	serializer, err := articles.NewArticlesService().Get(page)
	if err != nil {
		serializer.ErrorJson(w, err)
	}

	serializer.Encode(w, http.StatusAccepted)
}
