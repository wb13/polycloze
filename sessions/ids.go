// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package sessions

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"
)

// Generates a cryptographically secure random 128-bit string in base64.
// Results aren't guaranteed to be unique; use `generateUniqueID` instead.
func generateID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

// Generates a session ID that's not already in use.
func generateUniqueID(db *sql.DB) (string, error) {
	for {
		id, err := generateID()
		if err != nil {
			return "", fmt.Errorf("failed to generate a unique ID: %v", err)
		}

		err = reserveID(db, id)
		if err == nil {
			return id, nil
		}

		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			continue
		}
		return "", fmt.Errorf("failed to generate a unique ID: %v", err)
	}
}

// Reserves session ID so others can't use it.
// Returns an error if the ID is already in use.
func reserveID(db *sql.DB, id string) error {
	query := `INSERT INTO user_session (session_id) VALUES (?)`
	_, err := db.Exec(query, id)
	return err
}

// Deletes session ID from the database.
func deleteID(db *sql.DB, id string) error {
	query := `DELETE FROM user_session WHERE session_id = ?`
	_, err := db.Exec(query, id)
	return err
}
