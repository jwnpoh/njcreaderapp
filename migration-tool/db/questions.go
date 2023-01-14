// Package db provides functions and types relevant to the backend database for the article feed.
package db

import (
	"context"
	"fmt"
)

type Question struct {
	Year    string
	Number  string
	Wording string
}

type QuestionsDB map[string]Question

// MapQuestions maps a list of questions in a file named by filename and maps them to a questions database.
func InitQuestionsDB() (QuestionsDB, error) {
	qnDB := make(map[string]Question)

	// filename := "db/files/pastyressayqns.txt"
	// file, err := os.Open(filename)
	// if err != nil {
	// 	return qnDB, fmt.Errorf("unable to open file %s - %w", filename, err)
	// }
	// defer file.Close()

	// scanner := bufio.NewScanner(file)

	// for scanner.Scan() {
	// 	s := scanner.Text()
	// 	xs := strings.SplitN(s, " ", 3)

	// 	year := xs[0]
	// 	number := xs[1]
	// 	wording := xs[2]

	// 	qn := Question{year, number, wording}

	// 	key := year + " " + number

	//        _, ok := qnDB[key]; if ok {
	//            fmt.Printf("Duplicate found for %s %s\n", key, wording)
	//            continue
	//        }
	//        qnDB[key] = qn
	// }

	// if err := scanner.Err(); err != nil {
	// 	return qnDB, fmt.Errorf("problem scanning lines in file and mapping questions")
	// }

	ctx := context.Background()
	srv, err := newSheetsService(ctx)
	if err != nil {
		return qnDB, fmt.Errorf("unable to start Sheets service: %w", err)
	}

	var sd sheetData
	sd.ID = "1nY3sFjXXonSL43C3vPpfnEZO5b4SBXVSfhWdJkUzJS4" // DB new
	sd.Range = "Questions"

	resp, err := srv.Spreadsheets.Values.Get(sd.ID, sd.Range).Do()
	if err != nil {
		return qnDB, fmt.Errorf("unable to retrieve data from sheet: %v", err)
	}

	sd.Values = resp.Values

	if len(sd.Values) == 0 {
		return qnDB, fmt.Errorf("no data found")
	}

	for _, row := range sd.Values {
		year := fmt.Sprintf("%v", row[1])
		number := fmt.Sprintf("%v", row[2])
		wording := fmt.Sprintf("%v", row[3])

		qn := Question{year, number, wording}
		key := fmt.Sprintf("%v", row[0])

		_, ok := qnDB[key]
		if ok {
			fmt.Printf("Duplicate found for %s %s\n", key, wording)
			continue
		}
		qnDB[key] = qn
		// fmt.Println("Questions: ", qn)
	}

	return qnDB, nil
}
