// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

// auth-related handlers.
package api

import (
	"net/http"

	"github.com/lggruspe/polycloze/auth"
	"github.com/lggruspe/polycloze/sessions"
)

func hasUserID(data map[string]any) bool {
	val, ok := data["userID"]
	if !ok {
		return false
	}
	switch val := val.(type) {
	case int:
		return val >= 0
	default:
		return false
	}
}

func hasUsername(data map[string]any) bool {
	val, ok := data["username"]
	if !ok {
		return false
	}
	switch val := val.(type) {
	case string:
		return len(val) > 0
	default:
		return false
	}
}

// Checks if user is authenticated.
func isSignedIn(s *sessions.Session) bool {
	return hasUserID(s.Data) && hasUsername(s.Data)
}

// HandlerFunc for user registrations.
func handleRegister(w http.ResponseWriter, r *http.Request) {
	// Redirect to home page if already signed in.
	db := auth.GetDB(r)
	s, err := sessions.ResumeSession(db, w, r)
	if err == nil && isSignedIn(s) {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data := make(map[string]any)
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if auth.Register(db, username, password) == nil {
			http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
			return
		}
		data["message"] = "This username is unavailable. Try another one."
	}

	if err := renderTemplate(w, "register.html", data); err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
	}
}

// HandlerFunc for signing in.
func handleSignIn(w http.ResponseWriter, r *http.Request) {
	db := auth.GetDB(r)
	s, err := sessions.ResumeSession(db, w, r)
	if err == nil && isSignedIn(s) {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data := make(map[string]any)
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		userID, err := auth.Authenticate(db, username, password)
		if err != nil {
			data["message"] = "Incorrect username or password."
			goto fail
		}

		s, err = sessions.StartSession(db, w, r)
		if err != nil {
			data["message"] = "Authentication failed."
			goto fail
		}

		s.Data["userID"] = userID
		s.Data["username"] = username
		if sessions.SaveData(db, s) != nil {
			data["message"] = "Authentication failed."
			goto fail
		}

		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

fail:
	if err := renderTemplate(w, "signin.html", data); err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
	}
}

// HandlerFunc for signing out.
func handleSignOut(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}

	db := auth.GetDB(r)
	if err := sessions.EndSession(db, w, r); err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
