package articles

type Article struct {
	Title     string     `json:"title"`
	URL       string     `json:"url"`
	Topics    []string   `json:"topics"`
	Questions []Question `json:"questions"`
	DateAdded string     `json:"dateAdded"`
	UnixTime  int64      `json:"unixTime,omitempty"`
}

type ArticleSeries struct {
	Series []Article `json:"data"`
}

type Question struct {
	Year    string
	Number  string
	Wording string
}
