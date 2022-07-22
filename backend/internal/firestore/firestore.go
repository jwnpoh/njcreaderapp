package fire

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

type FireStoreRepo struct {
	Ctx    context.Context
	Client *firestore.Client
}

const projectID = "the-njc-reader-244fd"

func NewFireStoreRepo() (*FireStoreRepo, error) {
	godotenv.Load(".env")
	credentials := os.Getenv("CREDENTIALS")
	ctx := context.Background()
	sa := option.WithCredentialsJSON([]byte(credentials))

	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: projectID}, sa)
	if err != nil {
		return nil, errors.New("Unable to initalize Firestore.")
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to initalize Firestore client - %v", err))
	}

	// return client, nil
	repo := FireStoreRepo{
		Ctx:    ctx,
		Client: client,
	}

	return &repo, nil
}

func (f *FireStoreRepo) GetAll() (*firestore.DocumentIterator, error) {
	iter := f.Client.Collection("articles").Documents(f.Ctx)
	return iter, nil
}

func (f *FireStoreRepo) Search(term string) (*firestore.DocumentIterator, error) {
	defer f.Client.Close()

	return nil, nil
}

func (f *FireStoreRepo) Store(articles []map[string]interface{}) error {

	for _, article := range articles {
		_, err := f.Client.Collection("articles").Doc(fmt.Sprintf("%d", article["unixtime"])).Set(f.Ctx, article)
		if err != nil {
			log.Println("Could not add article", article["Title"])
		}
	}

	return nil
}
