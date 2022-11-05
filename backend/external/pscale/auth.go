package pscale

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
)

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

func (ps *AuthDB) InsertToken(token *core.Token) error {
	query := "INSERT INTO sessions (token, userID, expiry, hash) VALUES (?, ?, ?, ?)"

	tx, err := ps.DB.Begin()
	if err != nil {
		return fmt.Errorf("PScaleArticles: unable to begin tx for inserting session token to db - %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(query, token.PlainToken, token.UserID, token.Expiry, token.Hash)
	if err != nil {
		return fmt.Errorf("PScaleArticles: unable to insert session token for %d to db - %w", token.UserID, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("PScaleArticles: unable to commit tx to db - %w", err)
	}

	return nil
}

func (ps *AuthDB) GetToken(tokenHash string) (*core.Token, error) {
	query := "SELECT token, userID, expiry, hash FROM sessions WHERE hash = ?"

	row := ps.DB.QueryRowx(query, tokenHash)

	var plainToken, hash, expiryString string
	var userID int
	var expiry time.Time

	err := row.Scan(&plainToken, &userID, &expiryString, &hash)
	if err != nil {
		return nil, fmt.Errorf("error scanning row to get user- %w", err)
	}

	expiry, err = time.Parse("2006-01-02 15:04:05", expiryString)
	if err != nil {
		return nil, fmt.Errorf("error parsing expiry- %w", err)
	}

	token := &core.Token{
		PlainToken: plainToken,
		UserID:     userID,
		Expiry:     expiry,
		Hash:       hash,
	}

	return token, nil
}

func (ps *AuthDB) DeleteToken(userID int) error {
	query := "DELETE FROM sessions WHERE userID = ?"

	tx, err := ps.DB.Begin()
	if err != nil {
		return fmt.Errorf("unable to begin db tx for deleting session token for userID %d - %w", userID, err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("unable to delete session token for userID %d from database - %w", userID, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("unable to commit tx to db for updating user - %w", err)
	}

	return nil
}
