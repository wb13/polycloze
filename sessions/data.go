// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package sessions

import (
	"database/sql"
)

// Gets session data from the database.
// Returns an empty map even if the ID does not exist.
func getData(db *sql.DB, id string) map[string]any {
	var userID sql.NullInt32
	var username sql.NullString

	data := make(map[string]any)
	query := `SELECT user_id, username FROM user_session WHERE session_id = ?`
	if err := db.QueryRow(query, id).Scan(&userID, &username); err == nil {
		if userID.Valid && username.Valid {
			data["userID"] = int(userID.Int32)
			data["username"] = username.String
		}
	}
	return data
}

// Saves session data.
// The session must exist already.
// `SaveData` would still return `nil`, but wouldn't insert a new entry for the missing session.
func SaveData(db *sql.DB, s *Session) error {
	query := `UPDATE user_session SET user_id = ?, username = ?, updated = strftime('%s', 'now') WHERE session_id = ?`
	_, err := db.Exec(query, s.Data["userID"], s.Data["username"], s.ID)
	return err
}
