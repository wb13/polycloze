// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package sessions

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
)

// Name of cookie that stores session ID.
const cookieName = "id"

// Gets session cookie from client.
// Returns an error if no ID is found.
// Does not validate the cookie.
func getCookie(r *http.Request) (*http.Cookie, error) {
	return r.Cookie(cookieName)
}

// Checks if the session ID in the cookie is still valid (in DB and not expired).
func validateCookie(db *sql.DB, c *http.Cookie) error {
	if c.Name != cookieName {
		return errors.New("incorrect cookie name")
	}
	var id string
	// TODO scan expiry date for checking
	query := `SELECT session_id FROM user_session WHERE session_id = ?`
	if err := db.QueryRow(query, c.Value).Scan(&id); err != nil {
		return fmt.Errorf("invalid session ID: %w", err)
	}
	return nil
}

func setCookie(w http.ResponseWriter, id string) {
	c := http.Cookie{
		Name:     cookieName,
		Value:    id,
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   false,
	}
	http.SetCookie(w, &c)
}

func deleteCookie(w http.ResponseWriter) {
	c := http.Cookie{
		Name:     cookieName,
		Value:    "",
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   false,
		MaxAge:   -1,
	}
	http.SetCookie(w, &c)
}
