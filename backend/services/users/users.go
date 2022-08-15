package users

import (
	"fmt"

	"github.com/jwnpoh/njcreaderapp/backend/external/pscale"
	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
	"github.com/jwnpoh/njcreaderapp/backend/services/hasher"
)

type UserManager interface {
	InsertUser(*core.User) error
	GetUser(username string) (*core.User, error)
	DeleteUser(id int) error
	UpdateUserPassword(id int, newPasswordHash string) error
}

type userManager struct {
	db pscale.PScaleUsersDB
}

func NewUserManager() (UserManager, error) {
	db, err := pscale.NewUsersDB()
	if err != nil {
		return nil, fmt.Errorf("unable to initialize user manager service - %w", err)
	}
	return &userManager{db: db}, nil
}

func (um *userManager) InsertUser(user *core.User) error {
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

func (um *userManager) GetUser(email string) (*core.User, error) {
	user, err := um.db.GetUser(email)
	if err != nil {
		return nil, fmt.Errorf("userManager: unable to get user %s - %w", email, err)
	}

	return user, nil
}

func (um *userManager) UpdateUserPassword(id int, newPassword string) error {
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

func (um *userManager) DeleteUser(id int) error {
	err := um.db.DeleteUser(id)
	if err != nil {
		return fmt.Errorf("userManager: unable to delete user %d - %w", id, err)
	}

	return nil
}
