package hasher

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func GenerateHash(password string) ([]byte, error) {
	input := []byte(password)

	hashed, err := bcrypt.GenerateFromPassword(input, 10)
	if err != nil {
		return []byte(""), fmt.Errorf("hasher: unable to generate hash from password - %w", err)
	}

	return hashed, nil
}

func CheckHash(hash []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(password))
}
