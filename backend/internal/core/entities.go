package core

// Article is the main entity of the app.
type Article struct {
	ID              int      `json:"id"`
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
	ID       string `json:"id"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	Tags     string `json:"tags"`
	Date     string `json:"date"`
	MustRead string `json:"must_read"`
}

// Long is the entity that represents articles in the longer reads section
type Long struct {
	ID    int    `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
	URL   string `json:"url,omitempty"`
	Topic string `json:"topic,omitempty"`
}

// LongSeries is a slice of Longs to transport between database, appilcation, and frontend.
type LongSeries []Long

// LongPayload is a slice of Longs to be parsed into from the payload received from the frontend via the longer reads insertLong method
type LongPayload []Long

// LongTopics is a slice of strings representing all topics for long articles in the database.
type LongTopics []string

// Question is the entity representing past year questions.
type Question struct {
	Year    string
	Number  string
	Wording string
}

type KV struct {
	K string
	V int
}

type Stats struct {
	NumberofArticles         int  `json:"number_of_articles"`
	TopicsWithMostArticles   []KV `json:"topics_with_most_articles"`
	TopicsWithFewestArticles []KV `json:"topics_with_fewest_articles"`

	QuestionsWithMostArticles []KV `json:"questions_with_most_articles"`

	QuestionsWthFewestArticles []KV `json:"questions_with_fewest_articles"`
}

type Post struct {
	ID           int      `json:"id,omitempty"`
	UserID       int      `json:"user_id"`
	Author       string   `json:"author"`
	AuthorClass  string   `json:"author_class"`
	Likes        int      `json:"likes"`
	TLDR         string   `json:"tldr"`
	Examples     string   `json:"examples"`
	Notes        string   `json:"notes,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	Date         string   `json:"date"`
	Public       bool     `json:"public"`
	CreatedAt    int64    `json:"created_at,omitempty"`
	ArticleID    string   `json:"article_id"`
	ArticleTitle string   `json:"article_title"`
	ArticleURL   string   `json:"article_url"`
}

type PostPayload struct {
	UserID       int      `json:"user_id"`
	Likes        int      `json:"likes"`
	TLDR         string   `json:"tldr"`
	Examples     string   `json:"examples"`
	Notes        string   `json:"notes,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	Date         string   `json:"date"`
	Public       string   `json:"public"`
	ArticleID    string   `json:"article_id"`
	ArticleTitle string   `json:"article_title"`
	ArticleURL   string   `json:"article_url"`
}

type Posts []Post

type LikesList struct {
	PostID     int      `json:"post_id"`
	LikedByIDs []int    `json:"-"`
	LikedBy    []string `json:"liked_by"`
}

type Relations struct {
	UserID          int      `json:"user_id"`
	FollowingIDs    []int    `json:"following_ids"`
	FollowingUsers  []string `json:"following_users"`
	FollowedByIDs   []int    `json:"followed_by_ids"`
	FollowedByUsers []string `json:"followed_by_users"`
}
