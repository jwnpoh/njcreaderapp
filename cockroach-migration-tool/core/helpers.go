package core

import (
	"fmt"
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "1234567890"

func ParseUnixTime(date string) (int64, error) {
	t, err := time.Parse("Jan 02 2006", date)
	if err != nil {
		return 0, fmt.Errorf("unable to parse input date - %w", err)
	}
	res := t.Unix()
	return res, nil
}

func GenerateRandomString() string {
	rand.Seed(time.Now().UnixNano())

	res := make([]byte, 10)

	for i := range res {
		res[i] = charset[rand.Intn(len(charset))]
	}
	return string(res)
}
