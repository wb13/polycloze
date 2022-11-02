// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

// Per-session anti-CSRF tokens.
package sessions

import (
	"database/sql"
	"fmt"
)

// Returns a anti-CSRF token (connected to session in DB).
func CreateCSRFToken(db *sql.DB, sessionID string) (string, error) {
	token, err := generateID()
	if err != nil {
		return "", fmt.Errorf("failed to generate a CSRF token: %v", err)
	}

	query := `INSERT INTO csrf_token (session_id, token) VALUES (?, ?)`
	if _, err = db.Exec(query, sessionID, token); err != nil {
		return "", fmt.Errorf("failed to generate a CSRF token: %v", err)
	}
	return token, nil
}

// Deletes a CSRF token from the database.
func DeleteCSRFToken(db *sql.DB, sessionID, token string) error {
	query := `
		DELETE FROM csrf_token WHERE rowid in (
			SELECT min(rowid) AS rowid FROM csrf_token WHERE (session_id, token) = (?, ?)
		)
	`
	_, err := db.Exec(query, sessionID, token)
	return err
}

// Validates CSRF token.
func CheckCSRFToken(db *sql.DB, sessionID, token string) bool {
	var result string
	query := `SELECT token FROM csrf_token WHERE (session_id, token) = (?, ?) LIMIT 1`
	err := db.QueryRow(query, sessionID, token).Scan(&result)
	return err == nil
}
