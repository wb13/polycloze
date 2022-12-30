// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

// For sending messages to session users.
package sessions

import (
	"database/sql"
	"fmt"
	"time"
)

type Message struct {
	// Excludes sessionID field, because caller should already have it.
	Created time.Time
	Message string
	Kind    string
}

// Gets all recent messages to the user.
func getMessages(tx *sql.Tx, sessionID string) ([]Message, error) {
	query := `
		SELECT created, message, kind
		FROM message
		WHERE session_id = ?
		ORDER BY created ASC
	`

	rows, err := tx.Query(query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages from the database: %v", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var created int64
		var message Message
		if err := rows.Scan(&created, &message.Message, &message.Kind); err != nil {
			return nil, fmt.Errorf("failed to get messages from the database: %v", err)
		}
		message.Created = time.Unix(created, 0)
		messages = append(messages, message)
	}
	return messages, nil
}

// Deletes all messages to the user from the database.
func deleteMessages(tx *sql.Tx, sessionID string) error {
	query := `DELETE FROM message WHERE session_id = ?`
	if _, err := tx.Exec(query, sessionID); err != nil {
		return fmt.Errorf("failed to delete messages from the database: %v", err)
	}
	return nil
}

// Returns recent messages to user.
// Also deletes these messages from the db.
func (s *Session) Messages() ([]Message, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to get messages from the database: %v", err)
	}

	messages, err := getMessages(tx, s.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages from the database: %v", err)
	}

	if err := deleteMessages(tx, s.ID); err != nil {
		return nil, fmt.Errorf("failed to get messages from the database: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to get messages from the database: %v", err)
	}
	return messages, nil
}

// Saves message for user into the database.
func (s *Session) Message(kind string, message string) error {
	switch kind {
	case "success":
		break
	case "info":
		break
	case "warning":
		break
	case "error":
		break
	default:
		// Set `kind` to "error" if it has an invalid value.
		kind = "error"
	}

	query := `INSERT INTO message (session_id, message, kind) VALUES (?, ?, ?)`
	if _, err := s.db.Exec(query, s.ID, message, kind); err != nil {
		return fmt.Errorf("failed to save message for user: %v", err)
	}
	return nil
}

// Saves success message.
func (s *Session) SuccessMessage(message string) error {
	return s.Message("success", message)
}

// Saves info message.
func (s *Session) InfoMessage(message string) error {
	return s.Message("info", message)
}

// Saves warning message.
func (s *Session) WarningMessage(message string) error {
	return s.Message("warning", message)
}

// Saves error message.
func (s *Session) ErrorMessage(message string) error {
	return s.Message("error", message)
}
