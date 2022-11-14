package core

// Article is the main entity of the app.
type Article struct {
	ID              int      `json:"id,omitempty"`
	Title           string   `json:"title"`
	URL             string   `json:"url"`
	Topics          []string `json:"topics"`
	Questions       []string `json:"questions"`
	QuestionDisplay []string `json:"question_display"`
	Date            string   `json:"date"`
	MustRead        bool     `json:"must_read"`
	PublishedOn     int64    `json:"published_on,omitempty"`
}

// ArticleSeries is a slice of Article to transport between database, application, and frontend.
type ArticleSeries []Article

type ArticlePayload []struct {
	ID       string `json:"id,omitempty"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	Tags     string `json:"tags"`
	Date     string `json:"date"`
	MustRead string `json:"must_read"`
}

// Question is the entity representing past year questions.
type Question struct {
	Year    string
	Number  string
	Wording string
}
