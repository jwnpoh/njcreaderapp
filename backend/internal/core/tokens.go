package core

import "time"

type Token struct {
	PlainToken string    `json:"-"`
	UserID     int       `json:"-"`
	Expiry     time.Time `json:"expiry"`
	Hash       string    `json:"token"`
}
