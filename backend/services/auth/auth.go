package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
)

type AuthDB interface {
	InsertToken(token *core.Token) error
	RefreshToken(token *core.Token) error
	GetToken(tokenHash string) (*core.Token, error)
	DeleteToken(token string) error
	ClearTokens(userID int) error
}

type Authenticator struct {
	db AuthDB
}

func NewAuthenticator(authDB AuthDB) *Authenticator {
	return &Authenticator{authDB}
}

func (auth *Authenticator) CreateToken(userID int, timeToLife time.Duration) (*core.Token, error) {
	token := &core.Token{
		UserID: userID,
		Expiry: time.Now().Add(timeToLife),
	}

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.PlainToken = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.PlainToken))
	hashed := hash[:]
	hashString := hex.EncodeToString(hashed)
	token.Hash = hashString

	err = auth.db.InsertToken(token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (auth *Authenticator) RefreshToken(token *core.Token) error {
	token.Expiry = time.Now().Add(24 * time.Hour)

	err := auth.db.RefreshToken(token)
	if err != nil {
		return err
	}

	return nil
}

func (auth *Authenticator) DeleteToken(token string) error {
	err := auth.db.DeleteToken(token)
	if err != nil {
		return err
	}

	return nil
}

func (auth *Authenticator) AuthenticateToken(r *http.Request) (*core.Token, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("no authorization header received")
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, fmt.Errorf("no authorization header received")
	}

	token := headerParts[1]

	tok, err := auth.checkToken(token)
	if err != nil {
		return nil, err
	}

	return tok, nil
}

func (auth *Authenticator) checkToken(tokenString string) (*core.Token, error) {
	token, err := auth.db.GetToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("auth: no user session found for %s - %w", tokenString, err)
	}

	if time.Since(token.Expiry) > 0 {
		err = auth.DeleteToken(token.PlainToken)
		return nil, fmt.Errorf("auth: token has expired")
	}

	return token, nil
}
