package core

import "github.com/google/uuid"

type User struct {
	ID          uuid.UUID `json:"id,omitempty"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	Class       string    `json:"class"`
	LastLogin   string    `json:"last_login,omitempty"`
	DisplayName string    `json:"display_name"`
	Hash        []byte    `json:"-"`
}
