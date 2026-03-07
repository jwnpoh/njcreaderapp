package articles

// Place this file at: backend/services/articles/search_test.go
// Run with: go test ./services/articles/ -v

import "testing"

func TestSearchExact(t *testing.T) {
	summary(t)

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "two words joined with &",
			input: "climate change",
			want:  "climate & change",
		},
		{
			name:  "three words",
			input: "war and peace",
			want:  "war & and & peace",
		},
		{
			name:  "single word unchanged",
			input: "geopolitics",
			want:  "geopolitics",
		},
		{
			name:  "extra spaces between words are collapsed",
			input: "climate  change",
			want:  "climate & change",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := searchExact(tt.input)
			check(t, "result", got, tt.want)
		})
	}
}

func TestSearchAND(t *testing.T) {
	summary(t)

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "two terms",
			input: "climate AND change",
			want:  "climate & change",
		},
		{
			name:  "three terms",
			input: "climate AND change AND politics",
			want:  "climate & change & politics",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := searchAND(tt.input)
			check(t, "result", got, tt.want)
		})
	}
}

func TestSearchOR(t *testing.T) {
	summary(t)

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "two terms",
			input: "climate OR war",
			want:  "climate | war",
		},
		{
			name:  "three terms",
			input: "climate OR war OR peace",
			want:  "climate | war | peace",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := searchOR(tt.input)
			check(t, "result", got, tt.want)
		})
	}
}

func TestSearchNOT(t *testing.T) {
	summary(t)

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "two terms",
			input: "climate NOT war",
			want:  "climate & !war",
		},
		{
			name:  "no NOT keyword passes through unchanged",
			input: "climate war",
			want:  "climate war",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := searchNOT(tt.input)
			check(t, "result", got, tt.want)
		})
	}
}

func TestSearchQn(t *testing.T) {
	summary(t)

	t.Run("prepends isQn prefix", func(t *testing.T) {
		got := searchQn("2023 - Q5")
		check(t, "result", got, "isQn2023 - Q5")
	})
}
