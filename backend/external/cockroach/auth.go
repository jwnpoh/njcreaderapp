package cockroach

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
)

type AuthDB struct {
	DB *sqlx.DB
}

func NewAuthDB(dsn string) (*AuthDB, error) {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize cockroach database - %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("no response from cockroach  database - %w", err)
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)

	return &AuthDB{DB: db}, nil
}

func (ps *AuthDB) InsertToken(token *core.Token) error {
	query := "INSERT INTO sessions (token, userID, expiry, hash) VALUES ($1, $2, $3, $4)"

	tx, err := ps.DB.Begin()
	if err != nil {
		return fmt.Errorf("cockroach auth: unable to begin tx for inserting session token to db - %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(query, token.PlainToken, token.UserID, token.Expiry, token.Hash)
	if err != nil {
		return fmt.Errorf("cockroach auth: unable to insert session token for %d to db - %w", token.UserID, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("cockroach auth: unable to commit tx to db - %w", err)
	}

	return nil
}

func (ps *AuthDB) RefreshToken(token *core.Token) error {
	query := "UPDATE sessions SET expiry = $1 WHERE token = $2"

	tx, err := ps.DB.Begin()
	if err != nil {
		return fmt.Errorf("cockroach auth: unable to begin tx for refreshing session token to db - %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(query, token.Expiry, token.PlainToken)
	if err != nil {
		return fmt.Errorf("cockroach auth: unable to insert session token for %d to db - %w", token.UserID, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("cockroach auth: unable to commit tx to db - %w", err)
	}

	return nil
}

func (ps *AuthDB) GetToken(tokenString string) (*core.Token, error) {
	query := "SELECT token, userID, expiry, hash FROM sessions WHERE token = $1"

	row := ps.DB.QueryRowx(query, tokenString)

	var plainToken string
	var hash []byte
	var userID uuid.UUID
	var expiry time.Time

	err := row.Scan(&plainToken, &userID, &expiry, &hash)
	if err != nil {
		return nil, fmt.Errorf("error scanning row to get user- %w", err)
	}

	token := &core.Token{
		PlainToken: plainToken,
		UserID:     userID,
		Expiry:     expiry,
		Hash:       hash,
	}

	return token, nil
}

func (ps *AuthDB) DeleteToken(token string) error {
	query := "DELETE FROM sessions WHERE token = $1"

	tx, err := ps.DB.Begin()
	if err != nil {
		return fmt.Errorf("unable to begin db tx for deleting session token %s - %w", token, err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(query, token)
	if err != nil {
		return fmt.Errorf("unable to delete session token %s from database - %w", token, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("unable to commit tx to db for deleting session token - %w", err)
	}

	return nil
}

func (ps *AuthDB) ClearTokens(userID int) error {
	query := "DELETE FROM sessions WHERE userID = $1"

	tx, err := ps.DB.Begin()
	if err != nil {
		return fmt.Errorf("unable to begin db tx for clearing session tokens for userID %d - %w", userID, err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("unable to clear session tokens for userID %d from database - %w", userID, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("unable to commit tx to db for clearing session tokens - %w", err)
	}

	return nil
}
