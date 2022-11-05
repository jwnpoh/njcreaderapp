package articles

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	// "github.com/jwnpoh/njcreaderapp/backend/external/pscale"

	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
	"github.com/jwnpoh/njcreaderapp/backend/services/serializer"
)

type ArticleService interface {
	Get(offset int) (*core.ArticleSeries, error)
	Find(terms string) (*core.ArticleSeries, error)
	Store(data *core.ArticleSeries) error
}

type Articles struct {
	db ArticleService
}

// NewArticlesService returns an articleService object to implement methods to interact with PlanetScale database.
func NewArticlesService(articlesDB ArticleService) *Articles {
	return &Articles{db: articlesDB}
}

// Get gets up to 12 documents per page from PScale and serves them in pages of 12 articles each.
func (a *Articles) Get(page int) (serializer.Serializer, error) {
	// logic for pagination
	n := ((page - 1) * 10)
	if n < 0 {
		n = 0
	}

	// get articles from db
	articles, err := a.db.Get(n)
	if err != nil {
		return nil, fmt.Errorf("unable to get articles from db - %w", err)
	}

	return serializer.NewSerializer(false, fmt.Sprintf("got articles from page %d", n), articles), nil
}

// Find parses the query and sends it to database for querying results
func (a *Articles) Find(q string) (serializer.Serializer, error) {
	terms := checkQuery(q)
	fmt.Println("querying database with ", terms)

	articles, err := a.db.Find(terms)
	if err != nil {
		return serializer.NewSerializer(true, "no articles matched the query", articles), fmt.Errorf("unable to find articles relevant to the query %s from db - %w", terms, err)
	}

	return serializer.NewSerializer(false, fmt.Sprintf("found articles matching '%v'", terms), articles), nil
}

// Store parses the input time from front end admin dashboard to unix time, then sends the data to PlanetScale.
func (a *Articles) Store(data core.ArticleSeries) error {
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

func checkQuery(q string) string {
	searchYrAndQn := regexp.MustCompile(`^\d{4}\s?-?\s?(q|Q)?\d{1,2}$`)

	if searchYrAndQn.MatchString(q) {
		fmt.Println("found irregular question search", q)
		cutYear := regexp.MustCompile(`\d{4}`)
		year := cutYear.FindString(q)
		fmt.Println("year - ", year)
		q = strings.TrimLeft(q, year)
		cutQnNo := regexp.MustCompile(`(q|Q)?\d{1,2}`)
		qnNumber := strings.TrimLeft(strings.ToLower(cutQnNo.FindString(q)), "q")
		fmt.Println("q no - ", qnNumber)

		q = fmt.Sprintf("%s - Q%s", year, qnNumber)
		fmt.Println("reformatted question term ", q)
	}

	switch {
	case strings.Contains(q, " "):
		return searchExact(q)
	case strings.Contains(q, "AND"):
		return searchAND(q)
	case strings.Contains(q, "OR"):
		return searchOR(q)
	case strings.Contains(q, "NOT"):
		return searchNOT(q)
	default:
		return q
	}
}
