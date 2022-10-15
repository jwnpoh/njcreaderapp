package core

import "time"

type Token struct {
	Token  string    `json:"token"`
	UserID int       `json:"-"`
	Expiry time.Time `json:"expiry"`
}
