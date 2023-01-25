package pscale

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
)

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

func (uDB *UsersDB) InsertUsers(users *[]core.User) error {
	query := "INSERT INTO users (email, hash, role, last_login, display_name, class) VALUES (?, ?, ?, ?, ?, ?)"

	tx, err := uDB.DB.Begin()
	if err != nil {
		return fmt.Errorf("PScaleUsers: unable to begin tx for adding articles input to db - %w", err)
	}
	defer tx.Rollback()

	for i, user := range *users {
		_, err = tx.Exec(query, user.Email, user.Hash, user.Role, user.LastLogin, user.DisplayName, user.Class)
		if err != nil {
			return fmt.Errorf("PScaleUsers: unable to add user %s to db - %w", user.Email, err)
		}

		if i >= 0 && i%200 == 0 || i == len(*users)-1 {
			err = tx.Commit()
			if err != nil {
				return fmt.Errorf("PScaleUsers: unable to commit tx to db - %w", err)
			}
			tx, err = uDB.DB.Begin()
			if err != nil {
				return fmt.Errorf("PScaleUsers: unable to begin tx for adding articles input to db - %w", err)
			}
			defer tx.Rollback()
		}
	}

	return nil
}

func (uDB *UsersDB) GetUser(field string, value any) (*core.User, error) {
	query := fmt.Sprintf("SELECT * FROM users WHERE %s = ?", field)

	row := uDB.DB.QueryRowx(query, value)

	var email, hash, role, class, displayName, lastLogin string
	var id int
	err := row.Scan(&id, &email, &hash, &role, &lastLogin, &displayName, &class)
	if err != nil {
		return nil, fmt.Errorf("PScaleUsers: error scanning row to get user- %w", err)
	}

	user := core.User{
		ID:          id,
		Email:       email,
		Hash:        hash,
		DisplayName: displayName,
		Role:        role,
		Class:       class,
		LastLogin:   lastLogin,
	}

	return &user, nil
}

func (uDB *UsersDB) UpdateUser(newUser *core.User) error {
	query := "UPDATE users SET hash = ?, last_login = ?, display_name = ? WHERE id = ?"

	tx, err := uDB.DB.Begin()
	if err != nil {
		return fmt.Errorf("PScaleUsers: unable to begin db tx for updating user %d - %w", newUser.ID, err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(query, newUser.Hash, newUser.LastLogin, newUser.DisplayName, newUser.ID)
	if err != nil {
		return fmt.Errorf("PScaleUsers: unable to update user %d - %w", newUser.ID, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("PScaleUsers: unable to commit tx to db for updating user - %w", err)
	}

	return nil
}

func (uDB *UsersDB) DeleteUser(id int) error {
	query := "DELETE FROM users WHERE id = ?"

	tx, err := uDB.DB.Begin()
	if err != nil {
		return fmt.Errorf("PScaleUsers: unable to begin db tx for deleting user %d - %w", id, err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(query, id)
	if err != nil {
		return fmt.Errorf("PScaleUsers: unable to delete user %d - %w", id, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("PScaleUsers: unable to commit tx to db for updating user - %w", err)
	}

	return nil
}
