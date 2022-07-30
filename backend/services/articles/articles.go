package articles

import (
	"fmt"

	"github.com/jwnpoh/njcreaderapp/backend/pscale"

	// "github.com/jwnpoh/njcreaderapp/backend/internal/core"
	"github.com/jwnpoh/njcreaderapp/backend/serializer"
)

type articleService struct {
	db *pscale.PscaleDB
}

// NewArticlesService returns an articleService object to implement methods to interact with Firestore repo.
func NewArticlesService() (*articleService, error) {
	db, err := pscale.NewPscaleDB()
	if err != nil {
		return nil, fmt.Errorf("unable to initialize pscale database - %w", err)
	}
	return &articleService{db: db}, nil
}

// Get gets up to 12 documents per page from PScale and serves them in pages of 12 articles each.
func (a *articleService) Get(page int) (serializer.Serializer, error) {
	n := ((page - 1) * 12)
	if n < 0 {
		n = 0
	}

	articles, err := a.db.Get(n)
	if err != nil {
		return serializer.NewSerializer(true, fmt.Sprintf("unable to get articles"), nil), err
	}

	return serializer.NewSerializer(false, "hit the broker", articles), nil
}

// func (a *articleService) Store(data core.ArticleSeries) error {
// 	// go func indexer
// 	i := indexer.NewIndexer(data)
// 	index, err := i.Index()
// 	if err != nil {
// 		return err
// 	}
// 	// send index to firestore
// 	a.fireStoreRepo.Index(index)

// 	// send to Firestore
// 	articles := make(core.ArticleSeries, 0, 5)
// 	articles = append(articles, data...)
// 	err = a.fireStoreRepo.Store(&articles)
// 	if err != nil {
// 		return fmt.Errorf("unable to store articles - %w", err)
// 	}

// 	return nil
// }
