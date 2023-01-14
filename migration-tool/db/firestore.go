package db

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

const projectID = "the-njc-reader-244fd"
const pathToCredentialsFile = "./the-njc-reader-244fd-firebase-adminsdk-p1f8d-cc6c2575b4.json"

type fireStoreRepo struct {
	ctx    context.Context
	client *firestore.Client
}

func NewFireStoreApp(projectID string) (*fireStoreRepo, error) {
	log.Println("Initializing Firestore...")
	repo := fireStoreRepo{}
	ctx := context.Background()
	sa := option.WithCredentialsFile(pathToCredentialsFile)
	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: projectID}, sa)
	if err != nil {
		return nil, fmt.Errorf("Unable to initalize Firestore: %w", err)
	}

	log.Println("Initializing firestore client...")
	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("Unable to initalize Firestore client: %w", err)
	}

	repo.ctx = ctx
	repo.client = client
	return &repo, nil
}

func Add(database *ArticlesDBByDate, r *fireStoreRepo) error {
	defer r.client.Close()

	var count int
	index := make(map[string][]int64)

	for i, article := range *database {
		fmt.Printf("Added %v/%v articles\r", count, database.Len())
		article.Date += int64(i)
		// _, err := r.client.Collection("articles").Doc(fmt.Sprintf("%d", article.Date)).Set(r.ctx, map[string]interface{}{
		// 	"title":     article.Title,
		// 	"url":       article.URL,
		// 	"topics":    article.Topics,
		// 	"questions": article.Questions,
		// 	"date":      article.DisplayDate,
		// 	"unixtime":  article.Date,
		// })
		// if err != nil {
		// 	log.Println("Could not add article", article.Title)
		// }

		// index title
		nonWords := regexp.MustCompile(`\W`)
		title := article.Title
		title = strings.ToLower(title)
		titleTokens := strings.Fields(title)
	first:
		for _, v := range titleTokens {
			sanitized := nonWords.ReplaceAllString(v, "")
			for _, w := range index[sanitized] {
				if w == article.Date {
					// fmt.Println("Found value already in array for ", sanitized, "value is ", w)
					continue first
				}
			}
			index[sanitized] = append(index[sanitized], article.Date)
			// fmt.Println("key: ", sanitized, "value: ", index[sanitized])
		}

		//index questions
		questions := article.Questions
		for _, l := range questions {
			wording := l.Wording
			wording = strings.ToLower(wording)
			wordingTokens := strings.Fields(wording)
		second:
			for _, v := range wordingTokens {
				sanitized := nonWords.ReplaceAllString(v, "")
				for _, w := range index[sanitized] {
					if w == article.Date {
						// fmt.Println("Found value already in array for ", sanitized, "value is ", w)
						continue second
					}
				}
				index[sanitized] = append(index[sanitized], article.Date)
				// fmt.Println("key: ", sanitized, "value: ", index[sanitized])
			}
		}
		count++
	}

	for k, v := range index {
		_, err := r.client.Collection("index").Doc(k).Set(r.ctx, map[string]interface{}{
			"articles": v,
		})
		if err != nil {
			log.Println("Could not add index", k)
		}
	}

	return nil
}

func AddQuestions(qnDB QuestionsDB, r *fireStoreRepo) error {
	defer r.client.Close()

	var count int
	log.Println("Starting to add ", len(qnDB), " questions to Firestore...")

	for _, qn := range qnDB {
		fmt.Printf("Added %v/%v questions\r", count, len(qnDB))
		_, err := r.client.Collection("questions").Doc(qn.Year+qn.Number).Set(r.ctx, map[string]interface{}{
			"year":    qn.Year,
			"number":  qn.Number,
			"wording": qn.Wording,
		})

		if err != nil {
		}
		count++
	}

	return nil
}
