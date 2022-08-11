package pscale

import "github.com/jwnpoh/njcreaderapp/backend/internal/core"

type AuthDB interface {
	InsertToken(token *core.Token, user *core.User) error
	GetToken(tokenString string) (*core.User, error)
	DeleteToken(user *core.User) error
}
