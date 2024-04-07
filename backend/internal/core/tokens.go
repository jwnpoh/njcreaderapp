package core

import (
	"time"

	"github.com/google/uuid"
)

type Token struct {
	PlainToken string    `json:"token"`
	UserID     uuid.UUID `json:"-"`
	Expiry     time.Time `json:"expiry"`
	Hash       []byte    `json:"-"`
}
