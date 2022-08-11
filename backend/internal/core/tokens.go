package core

import "time"

type Token struct {
	Token  string    `json:"token"`
	UserID []byte    `json:"-"`
	Hash   []byte    `json:"-"`
	Expiry time.Time `json:"expiry"`
}
