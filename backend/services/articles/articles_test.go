package articles

// Place this file at: backend/services/articles/articles_test.go
// Run with: go test ./services/articles/ -v

import (
	"fmt"
	"reflect"
	"testing"
)

// =============================================================================
// Colour constants + test helpers
// =============================================================================

const (
	colGreen = "\033[32m"
	colRed   = "\033[31m"
	colBold  = "\033[1m"
	colReset = "\033[0m"
)

// check asserts equality and logs a coloured pass/fail line.
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

// summary registers a cleanup that prints a coloured PASS/FAIL banner
// after all subtests in a top-level test function have completed.
// Call it at the top of each Test* function.
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
// formatQuestionString
// =============================================================================

func TestFormatQuestionString(t *testing.T) {
	summary(t)

	tests := []struct {
		name       string
		input      string
		wantOutput string
		wantIsQn   bool
	}{
		// --- Should be classified as questions ---
		{
			name:       "standard format with space and hyphen",
			input:      "2023 - Q5",
			wantOutput: "2023 - Q5",
			wantIsQn:   true,
		},
		{
			name:       "no space around hyphen",
			input:      "2023-Q5",
			wantOutput: "2023 - Q5",
			wantIsQn:   true,
		},
		{
			name:       "lowercase q",
			input:      "2023 - q5",
			wantOutput: "2023 - Q5",
			wantIsQn:   true,
		},
		{
			name:       "two-digit question number",
			input:      "2019 - Q12",
			wantOutput: "2019 - Q12",
			wantIsQn:   true,
		},
		{
			name:       "extra whitespace padding",
			input:      "  2021 - Q3  ",
			wantOutput: "2021 - Q3",
			wantIsQn:   true,
		},
		{
			name:       "no Q prefix before number",
			input:      "2020-3",
			wantOutput: "2020 - Q3",
			wantIsQn:   true,
		},

		// --- Should NOT be classified as questions ---
		{
			name:       "plain topic word",
			input:      "geopolitics",
			wantOutput: "geopolitics",
			wantIsQn:   false,
		},
		{
			name:       "topic with spaces",
			input:      "work & life",
			wantOutput: "work & life",
			wantIsQn:   false,
		},
		{
			name:       "empty string",
			input:      "",
			wantOutput: "",
			wantIsQn:   false,
		},
		{
			name:       "year only, no question number",
			input:      "2023",
			wantOutput: "2023",
			wantIsQn:   false,
		},
		{
			name:       "topic that contains numbers but is not a question code",
			input:      "covid-19",
			wantOutput: "covid-19",
			wantIsQn:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOutput, gotIsQn := formatQuestionString(tt.input)
			check(t, "isQn", gotIsQn, tt.wantIsQn)
			check(t, "output", gotOutput, tt.wantOutput)
		})
	}
}

// =============================================================================
// splitTags
// =============================================================================

func TestSplitTags(t *testing.T) {
	summary(t)

	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "single topic",
			input: "geopolitics",
			want:  []string{"geopolitics"},
		},
		{
			name:  "two topics",
			input: "geopolitics;environment",
			want:  []string{"geopolitics", "environment"},
		},
		{
			name:  "topic and question code",
			input: "geopolitics;2023 - Q5",
			want:  []string{"geopolitics", "2023 - Q5"},
		},
		{
			name:  "trailing semicolon is stripped",
			input: "geopolitics;environment;",
			want:  []string{"geopolitics", "environment"},
		},
		{
			name:  "leading and trailing whitespace on the whole string",
			input: "  geopolitics;environment  ",
			want:  []string{"geopolitics", "environment"},
		},
		{
			name:  "empty string returns slice with one empty element",
			input: "",
			want:  []string{""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitTags(tt.input)
			check(t, "result", got, tt.want)
		})
	}
}

// =============================================================================
// checkQuery
// =============================================================================

func TestCheckQuery(t *testing.T) {
	summary(t)

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "question code is routed to searchQn",
			input: "2023 - Q5",
			want:  "isQn2023 - Q5",
		},
		{
			name:  "AND query",
			input: "climate AND change",
			want:  "climate & change",
		},
		{
			name:  "OR query",
			input: "climate OR war",
			want:  "climate | war",
		},
		{
			name:  "NOT query",
			input: "climate NOT war",
			want:  "climate & !war",
		},
		{
			name:  "multi-word query without operator uses exact",
			input: "climate change",
			want:  "climate & change",
		},
		{
			name:  "single word passes through unchanged",
			input: "geopolitics",
			want:  "geopolitics",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checkQuery(tt.input)
			check(t, "result", got, tt.want)
		})
	}
}
