package core

// Article is the main entity of the app.
type Article struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	URL         string   `json:"url"`
	Topics      []string `json:"topics"`
	Questions   []string `json:"questions"`
	Date        string   `json:"date"`
	PublishedOn int64    `json:"published_on,omitempty"`
}

// ArticleSeries is a slice of Article to transport between database, application, and frontend.
type ArticleSeries []Article

// Question is the entity representing past year questions.
type Question struct {
	Year    string
	Number  string
	Wording string
}
