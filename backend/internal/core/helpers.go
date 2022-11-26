package core

import (
	"fmt"
	"time"
)

func ParseUnixTime(date string) (int64, error) {
	t, err := time.Parse("Jan 02 2006", date)
	if err != nil {
		return 0, fmt.Errorf("unable to parse input date - %w", err)
	}
	res := t.Unix()
	return res, nil
}
