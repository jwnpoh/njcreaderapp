package main

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/jwnpoh/njcreaderapp/cockroach-migration-tool/core"
)

func migrateQuestionListTable(pscaleDB, cockroachDB *sqlx.DB) error {
	fmt.Println("getting questions from pscale...")
	questions, err := getPscaleQuestionList(pscaleDB)
	if err != nil {
		return fmt.Errorf("PScaleQuestions: unable to get question_list from pscale - %w", err)
	}
	fmt.Printf("got %d question_list from pscale\n", len(questions))

	fmt.Printf("attempting to insert %d questions from pscale to cockcroach...\n", len(questions))
	err = insertQuestionListToCockroach(cockroachDB, questions)
	if err != nil {
		return fmt.Errorf("CockroachQuestions: unable to insert question_list to cockroach - %w\n", err)
	}

	return nil
}

func getPscaleQuestionList(pscaleDB *sqlx.DB) ([]core.Question, error) {
	series := make([]core.Question, 0)

	query := "SELECT year, number, wording FROM question_list;"

	fmt.Println("running query")
	rows, err := pscaleDB.Queryx(query)
	if err != nil {
		return nil, fmt.Errorf("PScaleQuestions: unable to query question_list table - %w\n", err)
	}

	var count int
	for rows.Next() {
		var question core.Question
		err = rows.Scan(&question.Year, &question.Number, &question.Wording)
		if err != nil {
			return nil, fmt.Errorf("PScaleQuestions: error scanning row - %w\n", err)
		}

		series = append(series, question)
		count++
	}
	fmt.Printf("scanned a total of %d questions from pscale\n", len(series))

	return series, nil
}

func insertQuestionListToCockroach(cockroachDB *sqlx.DB, questions []core.Question) error {
	if len(questions) < 1 {
		fmt.Println("did not receive questions to insert to cockroach.")
		return fmt.Errorf("did not receive questions to insert to cockroach.")
	}

	query := "INSERT INTO question_list(question, year, number, wording) VALUES ($1, $2, $3, $4)"

	fmt.Println("running query")
	tx, err := cockroachDB.Begin()
	if err != nil {
		return fmt.Errorf("Cockroach: unable to begin tx for adding questions to cockroach - %w\n", err)
	}
	defer tx.Rollback()

	for i, question := range questions {
		qn := strings.Builder{}
		qn.WriteString(question.Year)
		qn.WriteString(" - Q")
		qn.WriteString(question.Number)
		fmt.Printf("inserting question %s\n", qn.String())
		_, err := tx.Exec(query, qn.String(), question.Year, question.Number, question.Wording)
		if err != nil {
			return fmt.Errorf("CockroachQuestions: unable to add question %s to cockroach\n", qn.String())
		}

		fmt.Printf("successfully inserted question #%d - %s\n", i+1, qn.String())
		if i >= 0 && i%200 == 0 || i == len(questions)-1 {
			err = tx.Commit()
			if err != nil {
				return fmt.Errorf("PScaleQuestions: unable to commit tx to db - %w", err)
			}
			tx, err = cockroachDB.Begin()
			if err != nil {
				return fmt.Errorf("PScaleQuestions: unable to begin tx for adding articles input to db - %w", err)
			}
			defer tx.Rollback()
		}
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("PScaleQuestions: unable to commit tx to db - %w", err)
	}

	return nil
}
