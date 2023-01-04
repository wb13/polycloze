// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

// Settings page.
package api

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"github.com/polycloze/polycloze/auth"
	"github.com/polycloze/polycloze/basedir"
	"github.com/polycloze/polycloze/database"
	"github.com/polycloze/polycloze/sessions"
)

func handleSettings(w http.ResponseWriter, r *http.Request) {
	db := auth.GetDB(r)
	s, err := sessions.ResumeSession(db, w, r)
	if err != nil || !s.IsSignedIn() {
		http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
		return
	}

	if r.Method == "POST" {
		username := s.Data["username"].(string)
		currentPassword := r.FormValue("current-password")
		newPassword := r.FormValue("new-password")
		csrfToken := r.FormValue("csrf-token")

		if !sessions.CheckCSRFToken(s.ID, csrfToken) {
			_ = s.ErrorMessage(
				"Something went wrong. Please try again.",
				"change-password",
			)
			goto fail
		}

		id, err := auth.Authenticate(db, username, currentPassword)
		if err != nil {
			log.Println(err)
			_ = s.ErrorMessage("Incorrect password.", "change-password")
			goto fail
		}

		if err := auth.ChangePassword(db, id, newPassword); err != nil {
			_ = s.ErrorMessage(
				"Something went wrong. Please try again.",
				"change-password",
			)
			goto fail
		}

		_ = s.SuccessMessage("Password updated.", "change-password")
	}

fail:
	// Get active course.
	userID := s.Data["userID"].(int)
	course, err := getUserActiveCourse(userID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	s.Data["course"] = course
	s.Data["csrfToken"] = sessions.CSRFToken(s.ID)
	s.Data["changePasswordMessages"], _ = s.Messages("change-password")
	s.Data["csvUploadMessages"], _ = s.Messages("csv-upload")
	s.Data["resetProgressMessages"], _ = s.Messages("reset-progress")
	renderTemplate(w, "settings.html", s.Data)
}

func handleResetProgress(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "expected POST request", http.StatusBadRequest)
		return
	}

	// Check if course exists.
	l1 := chi.URLParam(r, "l1")
	l2 := chi.URLParam(r, "l2")
	if !courseExists(l1, l2) {
		http.NotFound(w, r)
		return
	}

	// Check if user is signed in.
	db := auth.GetDB(r)
	s, err := sessions.ResumeSession(db, w, r)
	if err != nil || !s.IsSignedIn() {
		http.NotFound(w, r)
		return
	}
	userID := s.Data["userID"].(int)
	username := s.Data["username"].(string)
	csrfToken := r.FormValue("csrf-token")
	confirm := r.FormValue("confirm")

	// Check CSRF token.
	if !sessions.CheckCSRFToken(s.ID, csrfToken) {
		_ = s.ErrorMessage(
			"Something went wrong. Please try again.",
			"reset-progress",
		)
		goto fail
	}

	// Check confirmation string.
	if confirm != fmt.Sprintf("%v/%v-%v", username, l1, l2) {
		_ = s.ErrorMessage("Incorrect confirmation string.", "reset-progress")
		goto fail
	}

	if err := resetProgress(userID, l1, l2); err != nil {
		log.Println(err)
		_ = s.ErrorMessage(
			"Something went wrong. Please try again.",
			"reset-progress",
		)
		goto fail
	}

	_ = s.SuccessMessage("Progress has been reset.", "reset-progress")

fail:
	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}

// Resets course progress by deleting the review DB and re-initializing it.
func resetProgress(userID int, l1, l2 string) error {
	// TODO make this operation atomic

	// Delete review DB.
	path := basedir.Review(userID, l1, l2)
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to reset progress: %w", err)
	}

	// Re-initialize review DB.
	db, err := database.OpenUserDB(basedir.UserData(userID))
	if err != nil {
		return fmt.Errorf("failed to reset progress: %w", err)
	}
	defer db.Close()

	if err := setActiveCourse(db, userID, l1, l2); err != nil {
		return fmt.Errorf("failed to reset progress: %w", err)
	}
	return nil
}
