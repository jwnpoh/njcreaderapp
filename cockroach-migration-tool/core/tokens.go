package core

import (
	"time"

	"github.com/google/uuid"
)

type Token struct {
	PlainToken string    `json:"-"`
	UserID     uuid.UUID `json:"-"`
	Expiry     time.Time `json:"expiry"`
	Hash       string    `json:"token"`
}
