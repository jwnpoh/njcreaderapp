package fire

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

type FireStoreRepo struct {
	projectID   string
	credentials string
}

// NewFireStoreRepo initalizes a fireStoreRepo to interface with Firestore.
func NewFireStoreRepo() *FireStoreRepo {
	projectID := os.Getenv("PROJECT_ID")
	credentials := os.Getenv("CREDENTIALS")

	return &FireStoreRepo{
		projectID:   projectID,
		credentials: credentials,
	}
}

// Log stores a log item in firestore
func (r *FireStoreRepo) Log(v interface{}) error {
	// create Firestore app
	ctx := context.Background()
	sa := option.WithCredentialsJSON([]byte(r.credentials))

	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: r.projectID}, sa)
	if err != nil {
		return fmt.Errorf("unable to initalize Firestore - %w", err)
	}

	// create Firestore client
	client, err := app.Firestore(ctx)
	if err != nil {
		return fmt.Errorf("unable to initalize Firestore client - %v", err)
	}
	// remember to close the client
	defer client.Close()

	// add log document to firestore
	// t := strconv.Itoa(int(time.Now().Unix()))
	year := strconv.Itoa(time.Now().Year())
	month := time.Now().Month().String()
	date := time.Now().Format("2 Jan 2006 15:04")
	_, err = client.Collection("logs").Doc(year).Collection(month).Doc(date).Set(ctx, v)
	if err != nil {
		return fmt.Errorf("unable to add log entry to firestore - %w", err)
	}

	return nil
}
