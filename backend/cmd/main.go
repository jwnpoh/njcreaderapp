package main

import (
	"fmt"
	"log"

	"github.com/jwnpoh/njcreaderapp/backend/internal/articles-service"
	fire "github.com/jwnpoh/njcreaderapp/backend/internal/firestore"
)

func main() {
	service := articles.NewArticleService()
	repo, err := fire.NewFireStoreRepo()
	if err != nil {
		log.Fatal(err)
	}
	service.Repo = repo
	defer service.Repo.Client.Close()

	a1 := articles.Article{
		Title:     "Test article 1",
		URL:       "www.google.com",
		Topics:    []string{"Test", "article"},
		Questions: []articles.Question{},
		DateAdded: "Jun 20, 2020",
		UnixTime:  0,
	}

	articles := make([]articles.Article, 0)
	articles = append(articles, a1)

	service.Store(articles)

	fmt.Print(service.GetAll())
}
