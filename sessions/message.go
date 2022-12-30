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
	Context string // Empty means null/empty context.
}

// Checks if context args contains the empty string.
func containsNullContext(contexts []string) bool {
	for _, context := range contexts {
		if context == "" {
			return true
		}
	}
	return false
}

// Gets all recent messages to the user in the specified contexts.
func getMessages(
	tx *sql.Tx,
	sessionID string,
	contexts ...string,
) ([]Message, error) {
	if len(contexts) == 0 && !containsNullContext(contexts) {
		// Insert null context if empty.
		contexts = append(contexts, "")
	}
	query := `
		SELECT created, message, kind, context
		FROM message
		WHERE session_id = ? AND context = ?
		ORDER BY created ASC
	`

	var messages []Message
	for _, context := range contexts {
		rows, err := tx.Query(query, sessionID, context)
		if err != nil {
			return nil, fmt.Errorf("failed to get messages from the database: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var created int64
			var message Message
			err := rows.Scan(
				&created,
				&message.Message,
				&message.Kind,
				&message.Context,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to get messages from the database: %v", err)
			}
			message.Created = time.Unix(created, 0)
			messages = append(messages, message)
		}
	}
	return messages, nil
}

// Deletes all messages to the user from the database that are in the given
// contexts.
func deleteMessages(tx *sql.Tx, sessionID string, contexts ...string) error {
	if len(contexts) == 0 && !containsNullContext(contexts) {
		// Insert null context if empty.
		contexts = append(contexts, "")
	}

	query := `DELETE FROM message WHERE session_id = ? AND context = ?`
	for _, context := range contexts {
		if _, err := tx.Exec(query, sessionID, context); err != nil {
			return fmt.Errorf("failed to delete messages from the database: %v", err)
		}
	}
	return nil
}

// Returns recent messages to user.
// Also deletes these messages from the db.
// Only returns messages that belong in the specified contexts.
// If there are no context args, returns messages with a null context
// (those that were inserted with `Session.Message` without context args).
func (s *Session) Messages(contexts ...string) ([]Message, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to get messages from the database: %v", err)
	}

	messages, err := getMessages(tx, s.ID, contexts...)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages from the database: %v", err)
	}

	if err := deleteMessages(tx, s.ID, contexts...); err != nil {
		return nil, fmt.Errorf("failed to get messages from the database: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to get messages from the database: %v", err)
	}
	return messages, nil
}

// Saves message for user into the database.
// Adds one copy of the message for each context.
// If there are no context args, the message is added in a null context.
// See `Session.Messages` for info on how messages are retrieved from a
// specific context.
func (s *Session) Message(kind string, message string, contexts ...string) error {
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

	if len(contexts) == 0 && !containsNullContext(contexts) {
		// Insert null context if empty.
		contexts = append(contexts, "")
	}

	// Add one copy of the message for each context arg.
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to save message for user: %v", err)
	}

	query := `
		INSERT INTO message (session_id, message, kind, context)
		VALUES (?, ?, ?, ?)
	`
	for _, context := range contexts {
		if _, err := tx.Exec(query, s.ID, message, kind, context); err != nil {
			return fmt.Errorf("failed to save message for user: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to save message for user: %v", err)
	}
	return nil
}

// Saves success message.
// See Session.Message for info on context args.
func (s *Session) SuccessMessage(message string, contexts ...string) error {
	return s.Message("success", message, contexts...)
}

// Saves info message.
// See Session.Message for info on context args.
func (s *Session) InfoMessage(message string, contexts ...string) error {
	return s.Message("info", message, contexts...)
}

// Saves warning message.
// See Session.Message for info on context args.
func (s *Session) WarningMessage(message string, contexts ...string) error {
	return s.Message("warning", message, contexts...)
}

// Saves error message.
// See Session.Message for info on context args.
func (s *Session) ErrorMessage(message string, contexts ...string) error {
	return s.Message("error", message, contexts...)
}
