package users

import (
	"fmt"

	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
	"github.com/jwnpoh/njcreaderapp/backend/services/hasher"
)

type UsersDB interface {
	InsertUser(*core.User) error
	GetUser(field string, value any) (*core.User, error)
	DeleteUser(id int) error
	UpdateUser(id int, field, newValue string) error
}

type UserManager struct {
	db UsersDB
}

func NewUserManager(usersDB UsersDB) *UserManager {
	return &UserManager{db: usersDB}
}

func (um *UserManager) InsertUser(user *core.User) error {
	hash, err := hasher.GenerateHash(user.Hash)
	if err != nil {
		return fmt.Errorf("userManager: unable to generate hash from user input password - %w", err)
	}
	user.Hash = hash

	err = um.db.InsertUser(user)
	if err != nil {
		return fmt.Errorf("userManager: unable to insert new user - %w", err)
	}

	return nil
}

func (um *UserManager) GetUser(field string, value any) (*core.User, error) {
	user, err := um.db.GetUser(field, value)
	if err != nil {
		return nil, fmt.Errorf("userManager: unable to get user %s by %s - %w", value, field, err)
	}

	return user, nil
}

func (um *UserManager) UpdateUserPassword(id int, newPassword string) error {
	newPasswordHash, err := hasher.GenerateHash(newPassword)
	if err != nil {
		return fmt.Errorf("userManager: unable to generate hash from user input password - %w", err)
	}

	err = um.db.UpdateUser(id, "hash", newPasswordHash)
	if err != nil {
		return fmt.Errorf("userManager: unable to update user %d - %w", id, err)
	}

	return nil
}

func (um *UserManager) DeleteUser(id int) error {
	err := um.db.DeleteUser(id)
	if err != nil {
		return fmt.Errorf("userManager: unable to delete user %d - %w", id, err)
	}

	return nil
}
