// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

// User management.
package auth

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func saltHashPassword(password string) string {
	result, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(result)
}

func Register(db *sql.DB, username, password string) error {
	// NOTE allows empty string as password
	query := `INSERT INTO user (username, password) VALUES (?, ?)`
	hash := saltHashPassword(password)
	if _, err := db.Exec(query, username, hash); err != nil {
		return fmt.Errorf("unable to register user: %v", err)
	}
	return nil
}

// Validates credentials.
// Returns user ID on success.
func Authenticate(db *sql.DB, username, password string) (int, error) {
	var id int
	var hash string
	query := `SELECT id, password FROM user WHERE username = ?`
	err := db.QueryRow(query, username).Scan(&id, &hash)

	if err != nil && hash != "" {
		panic("something unexpected occurred")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return id, fmt.Errorf("unable to authenticate user: %v", err)
	}
	return id, nil
}
