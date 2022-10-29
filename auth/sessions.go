// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"log"
	"net/http"
)

const cookieName = "id"

type contextKey struct{}

type SessionData struct {
	UserID   int    // negative means none
	Username string // empty means none
}

type Session struct {
	ID   string
	Data SessionData

	db *sql.DB
}

// Generates 128-bit session ID in base64 encoding.
func generateSessionID() string {
	bytes := make([]byte, 16) // 128-bits
	if _, err := rand.Read(bytes); err != nil {
		panic("something unexpected occurred")
	}
	return base64.StdEncoding.EncodeToString(bytes)
}

// Get session data from database.
func getData(db *sql.DB, id string) (SessionData, error) {
	var data SessionData
	query := `SELECT user_id, username FROM user_session WHERE session_id = ?`
	err := db.QueryRow(query, id).Scan(&data.UserID, &data.Username)
	return data, err
}

// Deletes session data in database.
func deleteData(db *sql.DB, id string) error {
	query := `DELETE FROM user_session WHERE session_id = ?`
	_, err := db.Exec(query, id)
	return err
}

// Writes session data to database.
func saveData(db *sql.DB, session Session) error {
	query := `
		INSERT INTO user_session (session_id, user_id, username) VALUES (?, ?, ?)
			ON CONFLICT (session_id) DO UPDATE
			SET user_id = excluded.user_id, username = excluded.username
	`
	_, err := db.Exec(query, session.ID, session.Data.UserID, session.Data.Username)
	return err
}

// Creates cookie object with new default values.
func newCookie(name, value string) http.Cookie {
	return http.Cookie{
		Name:     name,
		Value:    value,
		SameSite: http.SameSiteStrictMode,
	}
}

func deleteCookie(name string) http.Cookie {
	return http.Cookie{
		Name:   name,
		Value:  "",
		MaxAge: -1,
	}
}

// Gets existing session from request context.
func GetSession(r *http.Request) Session {
	return r.Context().Value(contextKey{}).(Session)
}

// Gets existing session from cookie/db, or creates a new one.
func GenerateSession(db *sql.DB, r *http.Request) (Session, error) {
	session := Session{db: db}

	cookie, err := r.Cookie(cookieName)
	// NOTE should check cookie.Valid(), but incorrectly returns error when Expires is not set...
	if err == nil {
		if data, err := getData(db, cookie.Value); err == nil {
			session.ID = cookie.Value
			session.Data = data
			return session, nil
		}
	}
	session.ID = generateSessionID()
	session.Data.UserID = -1
	session.Data.Username = ""
	return session, nil
}

func (s Session) Save(w http.ResponseWriter) error {
	if err := saveData(s.db, s); err != nil {
		return err
	}
	cookie := newCookie(cookieName, s.ID)
	http.SetCookie(w, &cookie)
	return nil
}

// Deletes session data in cookies and database.
func (s Session) Delete(w http.ResponseWriter) error {
	if err := deleteData(s.db, s.ID); err != nil {
		return err
	}
	cookie := deleteCookie(cookieName)
	http.SetCookie(w, &cookie)
	return nil
}

// Gets session before each request.
// NOTE Doesn't auto-save sessions.
func Middleware(db *sql.DB) func(http.Handler) http.Handler {
	// Gets user session and stuffs it in the request context.
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := GenerateSession(db, r)
			if err != nil {
				log.Fatal(err)
			}

			ctx := context.WithValue(r.Context(), contextKey{}, session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
