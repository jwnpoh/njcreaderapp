package core

// Article is the main entity of the app.
type Article struct {
	Title     string     `json:"title"`
	URL       string     `json:"url"`
	Topics    []string   `json:"topics"`
	Questions []Question `json:"questions"`
	Date      string     `json:"date"`
	UnixTime  int64      `json:"unixTime,omitempty"`
}

// ArticleSeries is a slice of Article to transport between database, application, and frontend.
type ArticleSeries []Article


// Question is the entity representing past year questions.
type Question struct {
	Year    string
	Number  string
	Wording string
}
