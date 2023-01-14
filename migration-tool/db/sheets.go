// Package db provides functions and types relevant to the backend database for the article feed.
package db

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	// "golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "db/files/token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func newSheetsService(ctx context.Context) (*sheets.Service, error) {
	b, err := ioutil.ReadFile("db/files/njc-gp-newsfeed-002395a717c9.json")
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %v", err)
	}

	// b := os.Getenv("CREDENTIALS")
	// If modifying these scopes, delete your previously saved token.json.
	// config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	// if err != nil {
	// 	return nil, fmt.Errorf("unable to parse client secret file to config: %v", err)
	// }
	// client := getClient(config)

	srv, err := sheets.NewService(ctx, option.WithCredentialsJSON([]byte(b)))
	// srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return srv, fmt.Errorf("unable to retrieve Sheets client: %v", err)
	}

	return srv, nil
}

type sheetData struct {
	ID     string
	Range  string
	Values [][]interface{}
}

func getsheetData(srv *sheets.Service) (*sheetData, error) {
	var sd sheetData

	// sd.ID = "1Na9O_kaUSjde-drLyHIZevG847XgzF3FbakgaePDHuw" // DB old
	// sd.Range = "feed"

	sd.ID = "1nY3sFjXXonSL43C3vPpfnEZO5b4SBXVSfhWdJkUzJS4" // DB new
	sd.Range = "Articles"

	resp, err := srv.Spreadsheets.Values.Get(sd.ID, sd.Range).DateTimeRenderOption("FORMATTED_STRING").Do()
	if err != nil {
		return &sd, fmt.Errorf("unable to retrieve data from sheet: %v", err)
	}

	sd.Values = resp.Values

	return &sd, nil
}

// InitArticlesDB initialises the articles database taking data from the old articles DB. The data is to be sorted into the new data structures and migrated to Articles DB new.
func (database *ArticlesDBByDate) InitArticlesDB(ctx context.Context, qnDB QuestionsDB) error {
	srv, err := newSheetsService(ctx)
	if err != nil {
		return fmt.Errorf("unable to start Sheets service: %w", err)
	}

	data, err := getsheetData(srv)
	if err != nil {
		return fmt.Errorf("unable to get sheet data: %w", err)
	}

	if len(data.Values) == 0 {
		return fmt.Errorf("no data found")
	}

	var hasDuplicates bool
start:
	for _, row := range data.Values {
		a, err := NewArticle()
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		a.Title = fmt.Sprintf("%v", row[0])

		for _, article := range *database {
			if strings.ToLower(a.Title) == strings.ToLower(article.Title) {
				fmt.Printf("Duplicate entry: %s\n", a.Title)
				hasDuplicates = true
				continue start
			}
		}

		a.URL = fmt.Sprintf("%v", row[1])
		// a.SetDate(fmt.Sprintf("%v", row[4]))
		a.SetDate(fmt.Sprintf("%v", row[5]))

		// regex := regexp.MustCompile(`^\d{4}-Q\d{1,2}$`)
		// tags := row[5:]

		// for _, t := range tags {
		// 	tagString := fmt.Sprintf("%v", t)
		// 	if regex.MatchString(tagString) {
		// 		tagString = strings.ReplaceAll(tagString, "-Q", " ")
		// 		xt := strings.Split(tagString, " ")
		// 		year := xt[0]
		// 		number := xt[1]
		// 		a.SetQuestions(year, number, qnDB)
		// 	} else {
		// 		a.SetTopics(tagString)
		// 	}
		// }

		topics := fmt.Sprintf("%v", row[2])
		if topics == "" {
			fmt.Println(a.Title)
			continue
		}
		t := strings.Split(strings.ToLower(strings.ReplaceAll(topics, ",", "")), "\n")
		for _, topic := range t {
			switch topic {
			case "&":
				continue
			case "and":
				continue
			default:
				a.Topics = append(a.Topics, topic)
			}
		}

		questions := fmt.Sprintf("%v", row[3])
		if questions == "" {
			fmt.Println("Missing questions - " + a.Title)
			continue
		}
		q := strings.Split(questions, "\n")
		for _, r := range q {
			qn := qnDB[r]
			a.Questions = append(a.Questions, qn)
			// fmt.Printf("The questions for article %s are %v", a.Title, a.Questions)
		}
		*database = append(*database, *a)
	}
	sort.Sort(database)
	fmt.Println(strings.Repeat("-", 20))
	// sort.Sort(sort.Reverse(database))

	if hasDuplicates {
		return fmt.Errorf("duplicates in database. Please clean up first before proceeding with migration.")
	}
	return nil
}

// MigrateArticles backs up the articles database to a predefined, hard-coded Google Sheet.
func MigrateArticles(ctx context.Context, database *ArticlesDBByDate) error {
	// srv, err := newSheetsService(ctx)
	// if err != nil {
	// 	return fmt.Errorf("unable to start Sheets service: %w", err)
	// }

	// backupSheetID := "1nY3sFjXXonSL43C3vPpfnEZO5b4SBXVSfhWdJkUzJS4"
	// backupSheetName := "Articles"

	// var valueRange sheets.ValueRange
	// valueRange.Values = make([][]interface{}, 0, len(*database))

	// for _, j := range *database {
	// 	sTopics := strings.Builder{}
	// 	for i, k := range j.Topics {
	// 		if i == len(j.Topics)-1 {
	// 			sTopics.WriteString(string(k))
	// 			break
	// 		}
	// 		sTopics.WriteString(string(k) + "\n")
	// 	}
	// 	sQuestions := strings.Builder{}
	// 	sQuestionsKey := strings.Builder{}
	// 	for i, l := range j.Questions {
	// 		if i == len(j.Questions)-1 {
	// 			sQuestions.WriteString(fmt.Sprint(l.Year) + " " + fmt.Sprint(l.Number) + " " + l.Wording)
	// 			sQuestionsKey.WriteString(fmt.Sprint(l.Year) + " " + fmt.Sprint(l.Number))
	// 			break
	// 		}
	// 		sQuestions.WriteString(fmt.Sprint(l.Year) + " " + fmt.Sprint(l.Number) + " " + l.Wording + "\n")
	// 		sQuestionsKey.WriteString(fmt.Sprint(l.Year) + " " + fmt.Sprint(l.Number) + "\n")
	// 	}
	// 	record := make([]interface{}, 0, 6)
	// 	record = append(record, j.Title, j.URL, sTopics.String(), sQuestionsKey.String(), sQuestions.String(), j.DisplayDate)
	// 	valueRange.Values = append(valueRange.Values, record)
	// }

	// _, err = srv.Spreadsheets.Values.Update(backupSheetID, backupSheetName, &valueRange).ValueInputOption("RAW").Do()
	// if err != nil {
	// 	return fmt.Errorf("unable to backup data to backup sheet: %w", err)
	// }

	// repo, err := NewFireStoreApp(projectID)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = Add(database, repo)

	godotenv.Load(".env")
	db, err := newDB(os.Getenv("DSN"))
	if err != nil {
		log.Fatal(err)
	}

	err = PscaleAdd(database, db)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// MigrateQuestions backs up the questions database to a predefined, hard-coded Google Sheet.
func MigrateQuestions(ctx context.Context, qnDB QuestionsDB) error {
	// srv, err := newSheetsService(ctx)
	// if err != nil {
	// 	return fmt.Errorf("unable to start Sheets service: %w", err)
	// }

	// backupSheetID := "1nY3sFjXXonSL43C3vPpfnEZO5b4SBXVSfhWdJkUzJS4"
	// backupSheetName := "Questions"

	// var valueRange sheets.ValueRange
	// valueRange.Values = make([][]interface{}, 0, len(qnDB))

	// for k, v := range qnDB {
	// 	record := make([]interface{}, 0, 5)
	// 	record = append(record, k, v.Year, v.Number, v.Wording)
	// 	valueRange.Values = append(valueRange.Values, record)
	// }

	// _, err = srv.Spreadsheets.Values.Update(backupSheetID, backupSheetName, &valueRange).ValueInputOption("RAW").Do()
	// if err != nil {
	// 	return fmt.Errorf("unable to backup data to backup sheet: %w", err)
	// }

	// repo, err := NewFireStoreApp(projectID)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = AddQuestions(qnDB, repo)

	godotenv.Load(".env")
	db, err := newDB(os.Getenv("DSN"))
	if err != nil {
		log.Fatal(err)
	}

	err = PscaleAddQuestions(qnDB, db)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
