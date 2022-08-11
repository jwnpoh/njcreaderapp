package pscale

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
)

type PScaleUsersDB interface {
	InsertUser(*core.User) error
	GetUser(username string) (*core.User, error)
	DeleteUser(username, hash string) error
	UpdateUser(user *core.User, field, newValue string) error
}

func NewUsersDB() (PScaleUsersDB, error) {
	db, err := sqlx.Open("mysql", os.Getenv("PSCALE"))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize pscale database - %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("no response from pscale database - %w", err)
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)

	return &pscaleDB{DB: db}, nil
}
func (ps *pscaleDB) InsertUser(user *core.User) error {
	defer ps.DB.Close()

	query := "INSERT INTO users (email, hash, role, last_login) VALUES (?, ?, ?, ?)"

	tx, err := ps.DB.Begin()
	if err != nil {
		return fmt.Errorf("unable to begin tx for adding articles input to db - %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(query, user.Email, user.Hash, user.Role, user.LastLogin)
	if err != nil {
		return fmt.Errorf("unable to add user %s to db - %w", user.Email, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("unable to commit tx to db - %w", err)
	}

	return nil
}

func (ps *pscaleDB) GetUser(email string) (*core.User, error) {
	defer ps.DB.Close()

	query := "SELECT * FROM users WHERE email = ?"

	row := ps.DB.QueryRowx(query, email)

	var hash, role, lastLogin string
	var id int
	err := row.Scan(&id, &email, &hash, &role, &lastLogin)
	if err != nil {
		return nil, fmt.Errorf("error scanning row to get user- %w", err)
	}

	user := core.User{
		ID:        id,
		Email:     email,
		Hash:      hash,
		Role:      role,
		LastLogin: lastLogin,
	}

	return &user, nil
}

func (ps *pscaleDB) UpdateUser(user *core.User, field, newValue string) error {
	defer ps.DB.Close()

	if field != "hash" && field != "last_login" && field != "role" {
		return fmt.Errorf("field must be one of 'hash', 'last_login', or 'role'")
	}

	query := "UPDATE users SET ? = ? WHERE id = ? AND email = ?"

	tx, err := ps.DB.Begin()
	if err != nil {
		return fmt.Errorf("unable to begin db tx for updating user field %s - %w", field, err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(query, field, newValue, user.ID, user.Email)
	if err != nil {
		return fmt.Errorf("unable to update user field %s - %w", field, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("unable to commit tx to db for updating user - %w", err)
	}

	return nil
}

func (ps *pscaleDB) DeleteUser(username, hash string) error {
	return nil
}
func (ps *pscaleDB) InsertToken(token *core.Token, user *core.User) error {
	return nil
}
func (ps *pscaleDB) GetToken(tokenString string) (*core.User, error) {
	return &core.User{}, nil
}
func (ps *pscaleDB) DeleteToken(user *core.User) error {
	return nil
}
