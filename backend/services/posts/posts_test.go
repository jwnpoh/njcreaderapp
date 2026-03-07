package posts

// Place this file at: backend/services/posts/posts_test.go
// Run with: go test ./services/posts/ -v

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
)

// =============================================================================
// Colour helpers
// =============================================================================

const (
	colGreen = "\033[32m"
	colRed   = "\033[31m"
	colBold  = "\033[1m"
	colReset = "\033[0m"
)

func check(t *testing.T, name string, got, want any) {
	t.Helper()
	if reflect.DeepEqual(got, want) {
		t.Log(fmt.Sprintf("%s✓%s  %s", colGreen, colReset, name))
	} else {
		t.Errorf(
			"%s✗%s  %s%s%s\n     %sgot:%s  %#v\n     %swant:%s %#v",
			colRed, colReset,
			colBold, name, colReset,
			colRed, colReset, got,
			colGreen, colReset, want,
		)
	}
}

func summary(t *testing.T) {
	t.Helper()
	t.Cleanup(func() {
		if t.Failed() {
			t.Log(fmt.Sprintf("%s%s  FAIL  %s%s", colBold, colRed, t.Name(), colReset))
		} else {
			t.Log(fmt.Sprintf("%s%s  PASS  %s%s", colBold, colGreen, t.Name(), colReset))
		}
	})
}

// =============================================================================
// parsePostTags
//
// This is the most testable function in the posts service — pure input/output,
// no dependencies. It handles comma-splitting, whitespace trimming, hash
// stripping, lowercasing, and deduplication.
// =============================================================================

func TestParsePostTags(t *testing.T) {
	summary(t)

	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{
			name:  "single tag, no processing needed",
			input: []string{"geopolitics"},
			want:  []string{"geopolitics"},
		},
		{
			name:  "tags are lowercased",
			input: []string{"Geopolitics", "ENVIRONMENT"},
			want:  []string{"geopolitics", "environment"},
		},
		{
			name:  "hash prefix is stripped",
			input: []string{"#geopolitics", "#environment"},
			want:  []string{"geopolitics", "environment"},
		},
		{
			name:  "whitespace around tags is trimmed",
			input: []string{"  geopolitics  ", " environment "},
			want:  []string{"geopolitics", "environment"},
		},
		{
			name:  "comma-separated tags within a single element are split",
			input: []string{"geopolitics,environment,war"},
			want:  []string{"geopolitics", "environment", "war"},
		},
		{
			name:  "comma-separated tags mixed with standalone tags",
			input: []string{"geopolitics,environment", "war"},
			want:  []string{"geopolitics", "environment", "war"},
		},
		{
			name:  "duplicate tags are deduplicated",
			input: []string{"geopolitics", "geopolitics", "environment"},
			want:  []string{"geopolitics", "environment"},
		},
		{
			name:  "duplicate tags from comma split are also deduplicated",
			input: []string{"geopolitics,geopolitics"},
			want:  []string{"geopolitics"},
		},
		{
			name:  "hash, whitespace, and lowercase combined",
			input: []string{"#Geopolitics", "  #ENVIRONMENT  "},
			want:  []string{"geopolitics", "environment"},
		},
		{
			name:  "empty input returns empty slice",
			input: []string{},
			want:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parsePostTags(tt.input)
			check(t, "result", got, tt.want)
		})
	}
}

// =============================================================================
// parseNewPost
//
// Tests that parseNewPost correctly maps a PostPayload to a Post, with
// particular focus on the public flag ("on" → true, anything else → false)
// and that tags are processed through parsePostTags.
// =============================================================================

func TestParseNewPost(t *testing.T) {
	summary(t)

	userID := uuid.New()
	articleID := uuid.New()

	basePayload := func() *core.PostPayload {
		return &core.PostPayload{
			UserID:       userID,
			TLDR:         "A concise summary",
			Examples:     "Some examples",
			Notes:        "Additional notes",
			Tags:         []string{"geopolitics", "environment"},
			Public:       "on",
			ArticleID:    articleID,
			ArticleTitle: "Test Article",
			ArticleURL:   "https://example.com",
		}
	}

	t.Run("fields are mapped correctly from payload", func(t *testing.T) {
		post, err := parseNewPost(basePayload())

		check(t, "no error", err, nil)
		check(t, "TLDR", post.TLDR, "A concise summary")
		check(t, "Examples", post.Examples, "Some examples")
		check(t, "Notes", post.Notes, "Additional notes")
		check(t, "UserID", post.UserID, userID)
		check(t, "ArticleID", post.ArticleID, articleID)
		check(t, "ArticleTitle", post.ArticleTitle, "Test Article")
		check(t, "ArticleURL", post.ArticleURL, "https://example.com")
	})

	t.Run("public flag 'on' becomes true", func(t *testing.T) {
		payload := basePayload()
		payload.Public = "on"

		post, err := parseNewPost(payload)

		check(t, "no error", err, nil)
		check(t, "public is true", post.Public, true)
	})

	t.Run("public flag empty string becomes false", func(t *testing.T) {
		payload := basePayload()
		payload.Public = ""

		post, err := parseNewPost(payload)

		check(t, "no error", err, nil)
		check(t, "public is false", post.Public, false)
	})

	t.Run("public flag any other value becomes false", func(t *testing.T) {
		payload := basePayload()
		payload.Public = "off"

		post, err := parseNewPost(payload)

		check(t, "no error", err, nil)
		check(t, "public is false", post.Public, false)
	})

	t.Run("tags are processed through parsePostTags", func(t *testing.T) {
		payload := basePayload()
		// Mix of formats to confirm parsePostTags is being called
		payload.Tags = []string{"#Geopolitics", "  ENVIRONMENT  ", "geopolitics"}

		post, err := parseNewPost(payload)

		check(t, "no error", err, nil)
		// Expect lowercased, trimmed, hash-stripped, deduplicated
		check(t, "tags", post.Tags, []string{"geopolitics", "environment"})
	})

	t.Run("CreatedAt is set to a recent unix timestamp", func(t *testing.T) {
		post, err := parseNewPost(basePayload())

		check(t, "no error", err, nil)
		// CreatedAt should be a positive unix timestamp — we just verify
		// it's been set rather than pinning an exact value
		check(t, "CreatedAt is set", post.CreatedAt > 0, true)
	})
}
