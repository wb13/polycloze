// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package sessions

import (
	"database/sql"
	"fmt"
	"net/http"
)

type Session struct {
	ID   string
	Data map[string]any
}

// Starts a new session.
// Overwrites existing sessions, if any.
func StartSession(db *sql.DB, w http.ResponseWriter, r *http.Request) (*Session, error) {
	if err := EndSession(db, w, r); err != nil {
		return nil, fmt.Errorf("failed to start session: %v", err)
	}

	id, err := generateUniqueID(db)
	if err != nil {
		return nil, fmt.Errorf("failed to start session: %v", err)
	}

	setCookie(w, id)

	s := Session{
		ID:   id,
		Data: make(map[string]any),
	}
	return &s, nil
}

// Resumes an existing (valid) session.
// If there's none, returns an error.
func ResumeSession(db *sql.DB, w http.ResponseWriter, r *http.Request) (*Session, error) {
	c, err := getCookie(r)
	if err != nil {
		return nil, fmt.Errorf("failed to resume session: %v", err)
	}

	if err := validateCookie(db, c); err != nil {
		_ = EndSession(db, w, r)
		return nil, fmt.Errorf("failed to resume session: %v", err)
	}

	s := Session{
		ID:   c.Value,
		Data: getData(db, c.Value),
	}
	return &s, nil
}

// Resumes an existing (valid) session, or starts a new one if there's none yet.
func StartOrResumeSession(db *sql.DB, w http.ResponseWriter, r *http.Request) (*Session, error) {
	s, err := ResumeSession(db, w, r)
	if err == nil {
		return s, nil
	}
	return StartSession(db, w, r)
}

// Ends a session.
// Does nothing if there's no client session cookie.
func EndSession(db *sql.DB, w http.ResponseWriter, r *http.Request) error {
	var id string
	c, err := getCookie(r)
	if err == nil {
		id = c.Value
	}

	// Deleting empty ID doesn't delete any specific session, but deletes stale sessions.
	if err := deleteID(db, id); err != nil {
		return fmt.Errorf("failed to end session: %v", err)
	}

	// Deletes the cookie whether valid or not.
	deleteCookie(w)
	return nil
}
