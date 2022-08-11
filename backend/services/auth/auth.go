package auth

import (
	"time"

	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
)

type Authenticator interface {
	GetUser(email string) (*core.User, error)
	CreateToken(userID []byte, timeToLife time.Duration) (*core.Token, error)
	CheckToken(tokenString string) (*core.User, error)
}

type authenticator struct {
}

func NewAuthenticator() (Authenticator, error) {
	return &authenticator{}, nil
}

func (auth *authenticator) GetUser(email string) (*core.User, error) {

	return &core.User{}, nil
}
func (auth *authenticator) CreateToken(userID []byte, timeToLife time.Duration) (*core.Token, error) {
	return &core.Token{}, nil
}
func (auth *authenticator) CheckToken(tokenString string) (*core.User, error) {
	return &core.User{}, nil
}
