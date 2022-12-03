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
	return ProfanityCheck{
		Input:     input,
		IsProfane: goaway.IsProfane(input),
		Profanity: goaway.ExtractProfanity(input),
	}
}
