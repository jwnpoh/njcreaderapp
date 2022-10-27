package pscale

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
)

// type AuthDB interface {
// 	InsertToken(token *core.Token, user *core.User) error
// 	GetToken(tokenString string) (*core.Token, error)
// 	DeleteToken(user *core.Token) error
// }

type AuthDB struct {
	DB *sqlx.DB
}

func NewAuthDB(dsn string) (*AuthDB, error) {
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize pscale database - %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("no response from pscale database - %w", err)
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)

	return &AuthDB{DB: db}, nil
}

func (ps *AuthDB) InsertToken(token *core.Token, user *core.User) error {
	return nil
}

func (ps *AuthDB) GetToken(tokenString string) (*core.Token, error) {
	return &core.Token{}, nil
}

func (ps *AuthDB) DeleteToken(user *core.Token) error {
	return nil
}
