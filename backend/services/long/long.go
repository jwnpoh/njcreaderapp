package long

import (
	"bytes"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
	"github.com/jwnpoh/njcreaderapp/backend/services/serializer"
)

type LongsDB interface {
	GetTopics() (*core.LongTopics, error)
	Get(topic string) (*core.LongSeries, error)
	Store(data *core.LongPayload) error
	Update(data *core.Long) error
	Delete(ids string) error
}

type Longs struct {
	db LongsDB
}

// NewLongService returns a longService object to implement methods to interact with the database.
func NewLongService(longsDB LongsDB) *Longs {
	return &Longs{db: longsDB}
}

// GetTopics retrieves the full list of topics stored on the database.
func (l *Longs) GetTopics() (serializer.Serializer, error) {
	longTopics, err := l.db.GetTopics()
	if err != nil {
		return nil, fmt.Errorf("unable to get long topics from db - %w", err)
	}

	return serializer.NewSerializer(false, "got topics for long articles", longTopics), nil
}

// Get retrieves a slice of Long for a given topic to be sent to the client.
func (l *Longs) Get(topic string) (serializer.Serializer, error) {
	longSeries, err := l.db.Get(topic)
	if err != nil {
		return nil, fmt.Errorf("unable to get long articles of %s topic from db - %w", topic, err)
	}

	return serializer.NewSerializer(false, fmt.Sprintf("got long articles of %s topic", topic), longSeries), nil
}

// Store parses input from the client to be sent to the database.
func (l *Longs) Store(input string) (serializer.Serializer, error) {
	data, err := parseLongArticles(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse long articles input - %w", err)
	}

	err = l.db.Store(data)
	if err != nil {
		return nil, fmt.Errorf("unable to store long articles input - %w", err)
	}

	return serializer.NewSerializer(false, "successfully stored new long articles to db", nil), nil
}

// Update calls an update to a row in the long articles table in database, given an article id.
func (l *Longs) Update(data *core.Long) (serializer.Serializer, error) {
	err := l.db.Update(data)
	if err != nil {
		return nil, fmt.Errorf("unable to update long article %d - %w", data.ID, err)
	}

	return serializer.NewSerializer(false, fmt.Sprintf("successfully update long article %d", data.ID), nil), nil
}

// Delete takes an input of a slice of ids of string type to send to database for a bulk row delete query matching the article id.
func (l *Longs) Delete(input []string) (serializer.Serializer, error) {
	ids := strings.Join(input, ", ")

	err := l.db.Delete(ids)
	if err != nil {
		return nil, fmt.Errorf("unable to delete long articles %s - %w", ids, err)
	}

	return serializer.NewSerializer(false, "successfully deleted long articles", nil), nil
}

func parseLongArticles(input string) (*core.LongPayload, error) {
	input = strings.TrimSpace(input)
	rows := strings.Split(input, "\n")

	res := make(core.LongPayload, 0, len(rows))

	for _, row := range rows {
		var a core.Long

		row = strings.TrimSpace(row)
		r := strings.Split(row, ";")

		a.URL = strings.TrimSpace(r[0])
		a.Topic = strings.TrimSpace(r[1])

		title, err := getTitle(a.URL)
		if err != nil {
			return nil, err
		}

		a.Title = title

		res = append(res, a)
	}

	return &res, nil
}

func getTitle(url string) (string, error) {
	title := regexp.MustCompile(`\s?<meta[^p]*property=\"og:title\"\s?content=\"[^\"]*\"\s?\/?>`)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b2, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	t := title.Find(b2)

	regexHead := regexp.MustCompile(`\s?<meta[^p]*property=\"og:title\"\s?content=\"`)
	regexTail := regexp.MustCompile(`\"\s?\/?>`)

	head := regexHead.Find(t)
	tail := regexTail.Find(t)

	output := bytes.TrimPrefix(t, head)
	output = bytes.TrimSuffix(output, tail)
	output = bytes.TrimSpace(output)

	titleText := html.UnescapeString(string(output))
	if titleText == "" {
		return "", fmt.Errorf("could not find title - %w", err)
	}

	return titleText, nil
}
