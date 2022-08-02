package broker

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
	"github.com/jwnpoh/njcreaderapp/backend/serializer"
	"github.com/jwnpoh/njcreaderapp/backend/services/articles"
)

// GetLatest will be the "index" page of the application api.
func (b *broker) Get(w http.ResponseWriter, r *http.Request) {
	page := chi.URLParam(r, "page")

	n, err := strconv.Atoi(page)
	if err != nil {
		serializer.NewSerializer(true, "unable to parse request params", nil).ErrorJson(w, err)
	}

	service, err := articles.NewArticlesService()
	if err != nil {
		serializer.NewSerializer(true, "unable to start articles service", nil).ErrorJson(w, err)
	}

	serializer, err := service.Get(n)
	if err != nil {
		serializer.ErrorJson(w, err)
	}

	serializer.Encode(w, http.StatusAccepted)
}

func (b *broker) Store(w http.ResponseWriter, r *http.Request) {
	data := make(core.ArticleSeries, 0)

	err := serializer.Decode(w, r, &data)
	if err != nil {
		serializer.NewSerializer(true, "unable to decode input data", nil).ErrorJson(w, err)
	}

	service, err := articles.NewArticlesService()
	if err != nil {
		serializer.NewSerializer(true, "unable to start articles service", nil).ErrorJson(w, err)
	}

	err = service.Store(data)
	if err != nil {
		serializer.NewSerializer(true, "unable to store input data", nil).ErrorJson(w, err)
	}
}
