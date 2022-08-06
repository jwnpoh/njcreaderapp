package articles

import (
	"fmt"
	"time"

	"github.com/jwnpoh/njcreaderapp/backend/pscale"

	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
	"github.com/jwnpoh/njcreaderapp/backend/serializer"
)

type articleService struct {
	db pscale.PScale
}

// NewArticlesService returns an articleService object to implement methods to interact with PlanetScale database.
func NewArticlesService() (*articleService, error) {
	db, err := pscale.NewPscaleDB()
	if err != nil {
		return nil, fmt.Errorf("unable to initialize pscale database - %w", err)
	}
	return &articleService{db: db}, nil
}

// Get gets up to 12 documents per page from PScale and serves them in pages of 12 articles each.
func (a *articleService) Get(page int) (serializer.Serializer, error) {
	// logic for pagination
	n := ((page - 1) * 12)
	if n < 0 {
		n = 0
	}

	// get articles from db
	articles, err := a.db.Get(n)
	if err != nil {
		return nil, fmt.Errorf("unable to get articles from db - %w", err)
	}

	return serializer.NewSerializer(false, "hit the broker", articles), nil
}

// Find parses the query and sends it to database for querying results
func (a *articleService) Find(term string) (serializer.Serializer, error) {
	articles, err := a.db.Find(term)
	if err != nil {
		return serializer.NewSerializer(true, "no articles matched the query", articles), fmt.Errorf("unable to find articles relevant to the query %s from db - %w", term, err)
	}

	return serializer.NewSerializer(false, "hit the broker", articles), nil
}

// Store parses the input time from front end admin dashboard to unix time, then sends the data to PlanetScale.
func (a *articleService) Store(data core.ArticleSeries) error {
	// process string date to unixtime for PublishedOn field.
	for _, article := range data {
		t, err := time.Parse("Jan 2, 2006", article.Date)
		if err != nil {
			return fmt.Errorf("unable to parse input date - %w", err)
		}
		article.PublishedOn = t.Unix()
	}

	// send to PlanetScale
	err := a.db.Store(&data)
	if err != nil {
		return fmt.Errorf("unable to store articles - %w", err)
	}

	return nil
}
