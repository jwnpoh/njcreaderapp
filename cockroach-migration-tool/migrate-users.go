package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/jwnpoh/njcreaderapp/cockroach-migration-tool/core"
	"golang.org/x/crypto/bcrypt"
)

func migrateUsersTable(pscaleDB, cockroachDB *sqlx.DB) error {
	fmt.Println("getting users from pscale...")
	users, err := getPscaleUsers(pscaleDB)
	if err != nil {
		return fmt.Errorf("PScaleUsers: unable to get users from pscale - %w", err)
	}
	fmt.Printf("got %d users from pscale\n", len(users))

	fmt.Printf("attempting to insert %d users from pscale to cockcroach...\n", len(users))
	err = insertUsersToCockroach(cockroachDB, users)
	if err != nil {
		return fmt.Errorf("CockroachUsers: unable to insert users to cockroach - %w\n", err)
	}

	return nil
}

func getPscaleUsers(pscaleDB *sqlx.DB) ([]core.User, error) {
	series := make([]core.User, 0)

	query := "SELECT * FROM users;"

	fmt.Println("running query")
	rows, err := pscaleDB.Queryx(query)
	if err != nil {
		return nil, fmt.Errorf("PScaleUsers: unable to query users table - %w\n", err)
	}

	var count int
	for rows.Next() {
		var user core.User
		var id int
		err = rows.Scan(&id, &user.Email, &user.Hash, &user.Role, &user.LastLogin, &user.DisplayName, &user.Class)
		if err != nil {
			return nil, fmt.Errorf("PScaleUsers: error scanning row - %w\n", err)
		}

		series = append(series, user)
		count++
	}
	fmt.Printf("scanned a total of %d users from pscale\n", len(series))

	return series, nil
}

func insertUsersToCockroach(cockroachDB *sqlx.DB, users []core.User) error {
	if len(users) < 1 {
		fmt.Println("did not receive users to insert to cockroach.")
		return fmt.Errorf("did not receive usersto insert to cockroach.")
	}

	query := "INSERT INTO users(email, hash, role, last_login, display_name, class) VALUES ($1, $2, $3, $4, $5, $6)"

	fmt.Println("running query")
	tx, err := cockroachDB.Begin()
	if err != nil {
		return fmt.Errorf("Cockroach: unable to begin tx for adding users to cockroach - %w\n", err)
	}
	defer tx.Rollback()

	for i, user := range users {
		// hash, _ := GenerateHash("password")
		fmt.Printf("adding user #%d: \nEmail: %s\nHash: %s\nRole: %s\nLast login: %s\nDisplay name: %s\nClass: %s\n", i+1, user.Email, []byte(user.Hash), user.Role, user.LastLogin, user.DisplayName, user.Class)
		_, err := tx.Exec(query, user.Email, []byte(user.Hash), user.Role, user.LastLogin, user.DisplayName, user.Class)
		if err != nil {
			return fmt.Errorf("CockroachUsers: unable to add user %s to cockroach\n", user.DisplayName)
		}

		fmt.Printf("successfully inserted user #%d - %s\n", i+1, user.DisplayName)
		if i >= 0 && i%200 == 0 || i == len(users)-1 {
			err = tx.Commit()
			if err != nil {
				return fmt.Errorf("PScaleUsers: unable to commit tx to db - %w", err)
			}
			tx, err = cockroachDB.Begin()
			if err != nil {
				return fmt.Errorf("PScaleUsers: unable to begin tx for adding articles input to db - %w", err)
			}
			defer tx.Rollback()
		}
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("PScaleUsers: unable to commit tx to db - %w", err)
	}

	return nil
}

func GenerateHash(password string) (string, error) {
	input := []byte(password)

	hashed, err := bcrypt.GenerateFromPassword(input, 10)
	if err != nil {
		return "", fmt.Errorf("hasher: unable to generate hash from password - %w", err)
	}

	return string(hashed), nil
}
