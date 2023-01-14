// Package db provides functions and types relevant to the backend database for the article feed.
package db

import (
	"fmt"
	"sort"
	"strconv"
	// "strings"
	"time"
)

// Article is a struct representing a single entry in the articlesdb.
type Article struct {
	Title       string
	URL         string
	Topics      []string
	Questions   []Question
	DisplayDate string
	Date        int64
}

// Define topic-related data structures
type TopicsMap map[string]Article

// SetTopics is a wrapper around an append function to append multiple topics to the Article struct a.
func (a *Article) SetTopics(topics ...string) error {
	for _, j := range topics {
		a.Topics = append(a.Topics, j)
	}
	return nil
}

// SetQuestions sets the []Question for the given Article item.
func (a *Article) SetQuestions(year, number string, qnDB QuestionsDB) error {
	if _, err := strconv.Atoi(year); err != nil {
		return fmt.Errorf("the year input is not a number. try again")
	}
	if _, err := strconv.Atoi(number); err != nil {
		return fmt.Errorf("the question number input is not a number. try again")
	}
	key := year + " " + number
	qn := qnDB[key]
	a.Questions = append(a.Questions, qn)
	return nil
}

// SetDate is a wrapper around time.Parse to parse a date string to time.Time type, in order to call time.Unix() to return an int64 that makes the article sortable by date.
func (a *Article) SetDate(date string) error {
	// cleanDate := strings.ReplaceAll(date, ",", "")
	t, err := time.Parse("Jan 2, 2006", date)
	if err != nil {
		return fmt.Errorf("unable to parse date published - %w", err)
	}

	a.DisplayDate = date
	a.Date = t.Unix()
	return nil
}

// NewArticle returns an Article in order to populate fields for adding to the articles database.
func NewArticle() (*Article, error) {
	var a Article

	a.Topics = make([]string, 0)
	a.Questions = make([]Question, 0)

	return &a, nil
}

// ArticlesDBByDate is the database of all entries in the articlesdb. Entries are sorted in reverse order of date, with the most recent at index 0.
type ArticlesDBByDate []Article

// Implement sort.Sort interface
func (a ArticlesDBByDate) Len() int           { return len(a) }
func (a ArticlesDBByDate) Less(i, j int) bool { return a[i].Date < a[j].Date }
func (a ArticlesDBByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// NewArticlesDBByDate makes a slice of Articles to initialise the articles database.
func NewArticlesDBByDate() *ArticlesDBByDate {
	db := make(ArticlesDBByDate, 0, 50)

	return &db
}

// AddArticleToDB takes a pointer to an article that has already had all its fields populated and adds it to the articles database, sorting by most recent published date.
func (a *Article) AddArticleToDB(db *ArticlesDBByDate) error {
	*db = append(*db, *a)
	sort.Sort(sort.Reverse(db))

	return nil
}

// RemoveArticle is a function that the admin can invoke from the live app to remove any offending article.
func (db ArticlesDBByDate) RemoveArticle(index string) {
	j, _ := strconv.Atoi(index)
	copy(db[j:], db[j+1:])
	db[len(db)-1] = Article{}
}
