package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
)

type AuthDB interface {
	InsertToken(token *core.Token, user *core.User) error
	GetToken(tokenString string) (*core.Token, error)
	DeleteToken(user *core.Token) error
	// CreateToken(userID int, timeToLife time.Duration) (*core.Token, error)
	// AuthenticateToken(r *http.Request) (*core.User, error)
}

type Authenticator struct {
	db AuthDB
}

func NewAuthenticator(authDB AuthDB) *Authenticator {
	return &Authenticator{authDB}
}

func (auth *Authenticator) CreateToken(userID int, timeToLife time.Duration) (*core.Token, error) {
	// generate random string and insert to DB

	return &core.Token{}, nil
}

func (auth *Authenticator) checkToken(tokenString string) (*core.Token, error) {
	token, err := auth.db.GetToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("auth: no user session found %s - %w", tokenString, err)
	}

	return token, nil
}

func (auth *Authenticator) AuthenticateToken(r *http.Request) (int, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, fmt.Errorf("no authorization header received")
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return 0, fmt.Errorf("no authorization header received")
	}

	token := headerParts[1]
	if len(token) != 26 {
		return 0, fmt.Errorf("authentication token wrong size")
	}

	tok, err := auth.checkToken(token)
	if err != nil {
		return 0, fmt.Errorf("token not found")
	}

	return tok.UserID, nil
}
