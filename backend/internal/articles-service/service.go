package articles

import (
	"errors"
	"log"
	"time"

	fire "github.com/jwnpoh/njcreaderapp/backend/internal/firestore"
	"google.golang.org/api/iterator"
)

var (
	ErrNoArticlesFound       = errors.New("No articles found.")
	ErrArticleInfoIncomplete = errors.New("Article information is incomplete.")
	ErrArticleDateCorrupted  = errors.New("Article date added format error.")
)

type ArticleService struct {
	Series ArticleSeries
	Repo   *fire.FireStoreRepo
}

// GetAll implements Repository Interface
func (a *ArticleService) GetAll() (*ArticleSeries, error) {
	series := make([]Article, 0, 50)

	iter, err := a.Repo.GetAll()
	if err != nil {
		log.Printf("Unable to initialize Firestore to GetAll() - %v", err)
	}

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var article Article
		doc.DataTo(&article)
		series = append(series, article)
	}

	a.Series.Series = series
	return &a.Series, nil
}

// Search implements Repository Interface
func (a *ArticleService) Search(term string) (*ArticleSeries, error) {
	return &a.Series, nil
}

// Store implements Repository Interface
func (a *ArticleService) Store(aSeries []Article) error {
	data := make([]map[string]interface{}, 0, 5)
	for _, article := range aSeries {
		t, err := time.Parse("Jan 2, 2006", article.DateAdded)
		if err != nil {
			return ErrArticleDateCorrupted
		}
		article.UnixTime = t.Unix()

		a.Series.Series = append(a.Series.Series, article)

		m := make(map[string]interface{})
		m["title"] = article.Title
		m["url"] = article.URL
		m["date"] = article.DateAdded
		m["questions"] = article.Questions
		m["topics"] = article.Topics
		m["unixtime"] = article.UnixTime
		data = append(data, m)

		a.Repo.Store(data)
	}

	return nil
}

func NewArticleService() *ArticleService {
	return &ArticleService{}
}
