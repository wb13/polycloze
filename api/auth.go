// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

// auth-related handlers.
package api

import (
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/lggruspe/polycloze/auth"
	"github.com/lggruspe/polycloze/basedir"
	"github.com/lggruspe/polycloze/sessions"
)

// HandlerFunc for user registrations.
func handleRegister(w http.ResponseWriter, r *http.Request) {
	// Redirect to home page if already signed in.
	db := auth.GetDB(r)
	s, err := sessions.StartOrResumeSession(db, w, r)
	if err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	if s.IsSignedIn() {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")
		csrfToken := r.FormValue("csrf-token")

		if !sessions.CheckCSRFToken(s.ID, csrfToken) {
			_ = s.ErrorMessage("Something went wrong. Please try again.", "register")
			goto fail
		}
		if auth.Register(db, username, password) == nil {
			// `StatusTemporaryRedirect` also resends POST data to the next page.
			http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
			return
		}
		_ = s.ErrorMessage(
			"This username is unavailable. Try another one.",
			"register",
		)
	}

fail:
	messages, _ := s.Messages("register")
	data := map[string]any{
		"csrfToken": sessions.CSRFToken(s.ID),
		"messages":  messages,
	}
	renderTemplate(w, "register.html", data)
}

// HandlerFunc for signing in.
func handleSignIn(w http.ResponseWriter, r *http.Request) {
	db := auth.GetDB(r)
	s, err := sessions.StartOrResumeSession(db, w, r)
	if err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	var messages []sessions.Message
	if s.IsSignedIn() {
		goto success
	}

	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")
		csrfToken := r.FormValue("csrf-token")

		if !sessions.CheckCSRFToken(s.ID, csrfToken) {
			_ = s.ErrorMessage("Authentication failed.", "sign-in")
			goto fail
		}
		userID, err := auth.Authenticate(db, username, password)
		if err != nil {
			_ = s.ErrorMessage("Incorrect username or password.", "sign-in")
			goto fail
		}

		s.Data["userID"] = userID
		s.Data["username"] = username
		if sessions.SaveData(db, s) != nil {
			_ = s.ErrorMessage("Authentication failed.", "sign-in")
			goto fail
		}
		goto success
	}

fail:
	messages, _ = s.Messages("sign-in")
	renderTemplate(w, "signin.html", map[string]any{
		"csrfToken": sessions.CSRFToken(s.ID),
		"messages":  messages,
	})
	return

success:
	if err := initUserDirectory(s.Data["userID"].(int)); err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/welcome", http.StatusTemporaryRedirect)
}

// HandlerFunc for signing out.
func handleSignOut(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}

	db := auth.GetDB(r)
	s, err := sessions.ResumeSession(db, w, r)
	if err != nil || !s.IsSignedIn() {
		goto done
	}

	if err := sessions.EndSession(db, w, r); err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

done:
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func initUserDirectory(userID int) error {
	base := basedir.User(userID)
	reviews := path.Join(base, "reviews")

	if err := os.MkdirAll(base, 0o700); err != nil {
		return fmt.Errorf("failed to create user directory: %v", err)
	}
	if err := os.MkdirAll(reviews, 0o700); err != nil {
		return fmt.Errorf("failed to create user directory: %v", err)
	}
	return nil
}
