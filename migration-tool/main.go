package main

import (
	"context"
	"github.com/jwnpoh/njcreaderapp/migration-tool/db"
	"log"
)

func main() {
	ctx := context.Background()
	database := db.NewArticlesDBByDate()
	log.Println("Initialising questions DB...")
	qnDB, err := db.InitQuestionsDB()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Initialising articles DB...")
	if err := database.InitArticlesDB(ctx, qnDB); err != nil {
		log.Fatal(err)
	}

	log.Println("Migrating articles to planetscale...")
	if err := db.MigrateArticles(ctx, database); err != nil {
		log.Fatal(err)
	}

	log.Println("Migrating questions to pscale...")
	if err := db.MigrateQuestions(ctx, qnDB); err != nil {
		log.Fatal(err)
	}
}
