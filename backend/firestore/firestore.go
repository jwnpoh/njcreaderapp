package fire

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type fireStoreRepo struct {
	projectID   string
	credentials string
}

// NewFireStoreRepo initalizes a fireStoreRepo to interface with Firestore.
func NewFireStoreRepo() *fireStoreRepo {
	projectID := os.Getenv("PROJECT_ID")
	credentials := os.Getenv("CREDENTIALS")

	return &fireStoreRepo{
		projectID:   projectID,
		credentials: credentials,
	}
}

// Get queries Firestore for 144 articles.
func (r *fireStoreRepo) Get() (*core.ArticleSeries, error) {
	// create Firestore app
	ctx := context.Background()
	sa := option.WithCredentialsJSON([]byte(r.credentials))

	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: r.projectID}, sa)
	if err != nil {
		return nil, errors.New("Unable to initalize Firestore.")
	}

	// create Firestore client
	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to initalize Firestore client - %v", err))
	}
	// remember to close client
	defer client.Close()

	// parse data
	series := make(core.ArticleSeries, 0, 144)

	articles := client.Collection("articles")
	iter := articles.OrderBy("unixtime", firestore.Desc).Limit(144).Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var article core.Article
		doc.DataTo(&article)
		series = append(series, article)
	}

	return &series, nil
}

// Store takes a slice of Articles and creates new documents in Firestore.
func (r *fireStoreRepo) Store(articles *core.ArticleSeries) error {
	// create Firestore app
	ctx := context.Background()
	sa := option.WithCredentialsJSON([]byte(r.credentials))

	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: r.projectID}, sa)
	if err != nil {
		return errors.New("Unable to initalize Firestore.")
	}

	// create Firestore client
	client, err := app.Firestore(ctx)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to initalize Firestore client - %v", err))
	}
	// remember to close the client
	defer client.Close()

  // start adding documents to Firebase for each Article input
	for _, article := range *articles {
		_, err := client.Collection("articles").Doc(fmt.Sprintf("%d", article.UnixTime)).Set(ctx, map[string]interface{}{
			"title":     article.Title,
			"url":       article.URL,
			"topics":    article.Topics,
			"questions": article.Questions,
			"date":      article.Date,
			"unixtime":  article.UnixTime,
		})

		if err != nil {
			log.Println("Could not add article", article.Title)
		}
	}

	return nil
}
