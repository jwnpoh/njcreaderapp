package articles

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	// "github.com/jwnpoh/njcreaderapp/backend/external/pscale"

	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
	"github.com/jwnpoh/njcreaderapp/backend/services/serializer"
)

type ArticleService interface {
	Get(offset, limit int) (*core.ArticleSeries, error)
	Find(terms string) (*core.ArticleSeries, error)
	Store(data *core.ArticleSeries) error
	Update(data *core.ArticleSeries) error
	Delete(ids string) error
	GetQuestion(qnNo string) (string, error)
}

type Articles struct {
	db ArticleService
}

// NewArticlesService returns an articleService object to implement methods to interact with PlanetScale database.
func NewArticlesService(articlesDB ArticleService) *Articles {
	return &Articles{db: articlesDB}
}

// Get gets up to 10 documents per page from PScale and serves them in pages of 10 articles each.
func (a *Articles) Get(page, limit int) (serializer.Serializer, error) {
	// logic for pagination
	offset := ((page - 1) * limit)
	if offset < 0 {
		offset = 0
	}

	// get articles from db
	articles, err := a.db.Get(offset, limit)
	if err != nil {
		return nil, fmt.Errorf("unable to get articles from db - %w", err)
	}

	return serializer.NewSerializer(false, fmt.Sprintf("got articles from page %d", page), articles), nil
}

// Find parses the query and sends it to database for querying results
func (a *Articles) Find(q string) (serializer.Serializer, error) {
	terms := checkQuery(q)

	articles, err := a.db.Find(terms)
	if err != nil {
		return serializer.NewSerializer(true, "no articles matched the query", articles), fmt.Errorf("unable to find articles relevant to the query %s from db - %w", terms, err)
	}

	return serializer.NewSerializer(false, fmt.Sprintf("found articles matching '%v'", terms), articles), nil
}

// Store parses the input time from front end admin dashboard to unix time, then sends the data to PlanetScale.
func (a *Articles) Store(input core.ArticlePayload) error {
	data, err := a.parseNewArticles(input)
	if err != nil {
		return fmt.Errorf("unable to parse articles input - %w", err)
	}

	// send to PlanetScale
	err = a.db.Store(&data)
	if err != nil {
		return fmt.Errorf("unable to store articles - %w", err)
	}

	return nil
}

func (a *Articles) Update(input core.ArticlePayload) error {
	data, err := a.parseNewArticles(input)
	if err != nil {
		return fmt.Errorf("unable to parse articles input - %w", err)
	}

	// send to PlanetScale
	err = a.db.Update(&data)
	if err != nil {
		return fmt.Errorf("unable to store articles - %w", err)
	}

	return nil
}

func (a *Articles) Delete(input []string) error {
	ids := strings.Join(input, ", ")

	// send to Planetscale
	err := a.db.Delete(ids)
	if err != nil {
		return fmt.Errorf("unable to delete articles with id %s - %w", ids, err)
	}
	return nil
}

func checkQuery(q string) string {
	q, _ = formatQuestionString(q)

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

func formatQuestionString(q string) (string, bool) {
	searchYrAndQn := regexp.MustCompile(`\s?\d{4}(\s?)+-?(\s?)+(q|Q)?\d{1,2}`)
	var isQn bool

	if searchYrAndQn.MatchString(q) {
		q = strings.TrimSpace(q)
		cutYear := regexp.MustCompile(`\d{4}`)
		year := cutYear.FindString(q)
		q = strings.TrimLeft(q, year)
		cutQnNo := regexp.MustCompile(`(q|Q)?\d{1,2}`)
		qnNumber := strings.TrimLeft(strings.ToLower(cutQnNo.FindString(q)), "q")

		q = fmt.Sprintf("%s - Q%s", year, qnNumber)
		isQn = true
	}

	return q, isQn
}

func (a *Articles) parseNewArticles(input core.ArticlePayload) (core.ArticleSeries, error) {
	data := make(core.ArticleSeries, 0)

	for _, item := range input {
		date, err := parseUnixTime(item.Date)
		if err != nil {
			return nil, fmt.Errorf("unable to parse date input - %w", err)
		}

		tags := splitTags(item.Tags)

		questions, questionDisplay, topics, err := a.parseTags(tags)
		if err != nil {
			return nil, fmt.Errorf("unable to parse questions - %w", err)
		}

		var id int
		if item.ID != "" {
			id, err = strconv.Atoi(item.ID)
			if err != nil {
				return nil, fmt.Errorf("unable to id - %w", err)
			}
		}

		article := core.Article{
			ID:              id,
			Title:           item.Title,
			URL:             item.URL,
			Topics:          topics,
			Questions:       questions,
			QuestionDisplay: questionDisplay,
			PublishedOn:     date,
		}

		data = append(data, article)
	}

	return data, nil
}

func splitTags(tags string) []string {
	tags = strings.TrimSuffix(strings.TrimSpace(tags), ";")
	xtags := strings.Split(tags, ";")

	return xtags
}

func parseUnixTime(date string) (int64, error) {
	t, err := time.Parse("Jan 02 2006", date)
	if err != nil {
		return 0, fmt.Errorf("unable to parse input date - %w", err)
	}
	res := t.Unix()
	return res, nil
}

func (a *Articles) parseTags(tagsSlice []string) (questions, questionDisplay, topics []string, err error) {
	questions = make([]string, 0)
	questionDisplay = make([]string, 0)
	topics = make([]string, 0)

	for _, tag := range tagsSlice {
		qn, isQn := formatQuestionString(tag)
		if !isQn {
			topics = append(topics, strings.ToLower(tag))
			continue
		}
		questions = append(questions, qn)
		q, err := a.db.GetQuestion(qn)
		if err != nil {
			return questions, questionDisplay, topics, fmt.Errorf("unable to get question from database - %w", err)
		}
		questionDisplay = append(questionDisplay, q)
	}

	return questions, questionDisplay, topics, nil
}
