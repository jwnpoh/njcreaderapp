package profanity

import (
	"github.com/TwiN/go-away"
)

type ProfanityCheck struct {
	Input     string
	IsProfane bool
	Profanity string
}

func CheckProfanity(input string) ProfanityCheck {
	profanityDetector := goaway.NewProfanityDetector().WithSanitizeSpaces(false)

	return ProfanityCheck{
		Input:     input,
		IsProfane: profanityDetector.IsProfane(input),
		Profanity: profanityDetector.ExtractProfanity(input),
	}
}
