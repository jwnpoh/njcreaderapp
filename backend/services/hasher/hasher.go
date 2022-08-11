package hasher

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func GenerateHash(password string) (string, error) {
	input := []byte(password)

	hashed, err := bcrypt.GenerateFromPassword(input, 10)
	if err != nil {
		return "", fmt.Errorf("hasher: unable to generate hash from password - %w", err)
	}

	return string(hashed), nil
}

func CheckHash(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
