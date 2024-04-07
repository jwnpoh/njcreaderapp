package sheets

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/jwnpoh/njcreaderapp/backend/cmd/config"
	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SheetsDB struct {
	Service   *sheets.Service
	SheetID   string
	SheetName string
}

func NewSheetsService(ctx context.Context, cfg config.SheetsConfig) (*SheetsDB, error) {
	srv, err := sheets.NewService(ctx, option.WithCredentialsJSON([]byte(cfg.Credentials)))
	if err != nil {
		return nil, fmt.Errorf("unable to start new Sheets service - %w", err)
	}

	return &SheetsDB{
		Service:   srv,
		SheetID:   cfg.SheetID,
		SheetName: cfg.SheetName,
	}, nil
}

func (sheetsDB *SheetsDB) Get(offset, limit int) (*core.ArticleSeries, error) {
	return nil, nil
}

func (sheetsDB *SheetsDB) GetArticle(id int) (*core.Article, error) {
	return nil, nil

}

func (sheetsDB *SheetsDB) Find(terms string) (*core.ArticleSeries, error) {
	return nil, nil

}

func (sheetsDB *SheetsDB) Store(data *core.ArticleSeries) error {
	var valueRange sheets.ValueRange
	valueRange.Values = make([][]interface{}, 0, 1)

	for _, article := range *data {
		record := parseArticle(article)

		valueRange.Values = append(valueRange.Values, record)
	}
	_, err := sheetsDB.Service.Spreadsheets.Values.Append(sheetsDB.SheetID, sheetsDB.SheetName, &valueRange).InsertDataOption("INSERT_ROWS").ValueInputOption("RAW").Do()
	if err != nil {
		return fmt.Errorf("unable to append articles to sheets DB - %w", err)
	}

	return nil
}

func (sheetsDB *SheetsDB) Update(data *core.ArticleSeries) error {

	return nil
}

func (sheetsDB *SheetsDB) Delete(ids []string) error {

	return nil
}

func (sheetsDB *SheetsDB) GetQuestion(qnNo string) (string, error) {
	return "", nil
}

func parseArticle(article core.Article) []any {
	topics := parseTopics(article.Topics)

	keyRX := regexp.MustCompile(`\d{4}\s-\sQ\d{1,2}\s`)
	yearRX := regexp.MustCompile(`\d{4}`)
	qnRX := regexp.MustCompile(`Q\d{1,2}`)

	questionKeys := parseQuestionKeys(article.Questions, yearRX, qnRX)
	questionDisplay := parseQuestionDisplay(article.QuestionDisplay, yearRX, qnRX, keyRX)

	article.Date = time.Unix(article.PublishedOn, 0).Format("Jan 2, 2006")

	record := make([]interface{}, 0, 6)
	record = append(record, article.Title, article.URL, topics, questionKeys, questionDisplay, article.Date)

	return record
}

func parseTopics(topics []string) string {
	sTopics := strings.Builder{}
	for i, topic := range topics {
		if i == len(topics)-1 {
			sTopics.WriteString(strings.Title(topic))
			break
		}
		sTopics.WriteString(strings.Title(topic) + "\n")
	}

	return sTopics.String()
}

func parseQuestionKeys(questions []string, yearRX, qnRX *regexp.Regexp) string {
	sQuestionsKeys := strings.Builder{}

	for i, l := range questions {
		year := yearRX.FindString(l)
		qn := strings.ReplaceAll(qnRX.FindString(l), "Q", "")
		if i == len(questions)-1 {
			sQuestionsKeys.WriteString(year + " " + qn)
			break
		}
		sQuestionsKeys.WriteString(year + " " + qn + "\n")
	}
	return sQuestionsKeys.String()
}

func parseQuestionDisplay(questionDisplay []string, yearRX, qnRX, keyRX *regexp.Regexp) string {
	sQuestions := strings.Builder{}
	for i, l := range questionDisplay {
		year := yearRX.FindString(l)
		qn := strings.ReplaceAll(qnRX.FindString(l), "Q", "")
		wording := keyRX.ReplaceAllLiteralString(l, "")
		if i == len(questionDisplay)-1 {
			sQuestions.WriteString(year + " " + qn + " " + wording)
			break
		}
		sQuestions.WriteString(year + " " + qn + " " + wording + "\n")
	}
	return sQuestions.String()
}
