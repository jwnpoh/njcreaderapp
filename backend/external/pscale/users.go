package pscale

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
)

// type PScaleUsersDB interface {
// 	InsertUser(*core.User) error
// 	GetUser(field string, value any) (*core.User, error)
// 	DeleteUser(id int) error
// 	UpdateUser(id int, field, newValue string) error
// }

type UsersDB struct {
	DB *sqlx.DB
}

func NewUsersDB(dsn string) (*UsersDB, error) {
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

	return &UsersDB{DB: db}, nil
}

func (uDB *UsersDB) InsertUser(user *core.User) error {
	query := "INSERT INTO users (email, hash, role, last_login) VALUES (?, ?, ?, ?)"

	tx, err := uDB.DB.Begin()
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

func (uDB *UsersDB) GetUser(field string, value any) (*core.User, error) {
	query := fmt.Sprintf("SELECT * FROM users WHERE %s = ?", field)

	row := uDB.DB.QueryRowx(query, value)

	var email, hash, role, lastLogin string
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

func (uDB *UsersDB) UpdateUser(id int, field, newValue string) error {
	if field != "hash" && field != "last_login" && field != "role" {
		return fmt.Errorf("field must be one of 'hash', 'last_login', or 'role'")
	}

	query := fmt.Sprintf("UPDATE users SET %s = ? WHERE id = ?", field)

	tx, err := uDB.DB.Begin()
	if err != nil {
		return fmt.Errorf("unable to begin db tx for updating user field %s - %w", field, err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(query, newValue, id)
	if err != nil {
		return fmt.Errorf("unable to update user field %s - %w", field, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("unable to commit tx to db for updating user - %w", err)
	}

	return nil
}

func (uDB *UsersDB) DeleteUser(id int) error {
	query := "DELETE FROM users WHERE id = ?"

	tx, err := uDB.DB.Begin()
	if err != nil {
		return fmt.Errorf("unable to begin db tx for deleting user %d - %w", id, err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(query, id)
	if err != nil {
		return fmt.Errorf("unable to delete user %d - %w", id, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("unable to commit tx to db for updating user - %w", err)
	}

	return nil
}
