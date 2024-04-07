package users

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
	"github.com/jwnpoh/njcreaderapp/backend/services/hasher"
)

type UsersDB interface {
	InsertUsers(*[]core.User) error
	GetUser(field string, value any) (*core.User, error)
	DeleteUser(id uuid.UUID) error
	UpdateUser(*core.User) error
	UpdateClasses(*[]core.User) error
}

type UserManager struct {
	db UsersDB
}

func NewUserManager(usersDB UsersDB) *UserManager {
	return &UserManager{db: usersDB}
}

func (um *UserManager) InsertUsers(users *[]core.User) error {
	data := make([]core.User, 0, len(*users))

	for _, user := range *users {
		hash, err := hasher.GenerateHash("password")
		if err != nil {
			return fmt.Errorf("userManager: unable to generate hash from user input password - %w", err)
		}
		user.Hash = hash

		data = append(data, user)
	}

	err := um.db.InsertUsers(&data)
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

func (um *UserManager) UpdateUser(newUser *core.User, password string) error {
	newHash, err := hasher.GenerateHash(string(password))
	if err != nil {
		return fmt.Errorf("userManager: unable to generate hash from user input password - %w", err)
	}
	newUser.Hash = newHash

	err = um.db.UpdateUser(newUser)
	if err != nil {
		return fmt.Errorf("userManager: unable to update user %d - %w", newUser.ID, err)
	}

	return nil
}

func (um *UserManager) UpdateClasses(users *[]core.User) error {
	err := um.db.UpdateClasses(users)
	if err != nil {
		return fmt.Errorf("userManager: unable to insert new user - %w", err)
	}

	return nil
}

func (um *UserManager) DeleteUser(id uuid.UUID) error {
	err := um.db.DeleteUser(id)
	if err != nil {
		return fmt.Errorf("userManager: unable to delete user %d - %w", id, err)
	}

	return nil
}
