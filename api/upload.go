// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/polycloze/polycloze/auth"
	"github.com/polycloze/polycloze/basedir"
	"github.com/polycloze/polycloze/database"
	"github.com/polycloze/polycloze/replay"
	"github.com/polycloze/polycloze/sessions"
)

// Checks if uploaded file size is too big.
func isTooBig(size int64) bool {
	// Limit to 8MB.
	return size > 8*1024*1024
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
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
	var message string
	var success bool
	userID := s.Data["userID"].(int)

	// Check CSRF token.
	csrfToken := r.FormValue("csrf-token")
	if !sessions.CheckCSRFToken(s.ID, csrfToken) {
		http.Error(w, "Forbidden.", http.StatusForbidden)
		return
	}

	// Handle upload.
	file, header, err := r.FormFile("csv-upload")
	if err != nil {
		log.Println(err)
		message = "Something went wrong. Please try again."
		_ = s.ErrorMessage(message, "csv-upload")
		goto fail
	}

	if header.Header.Get("Content-Type") != "text/csv" {
		message = "Not a CSV file."
		_ = s.ErrorMessage(message, "csv-upload")
		goto fail
	}

	if isTooBig(header.Size) {
		message = "File is too big (>8MB)."
		_ = s.ErrorMessage(message, "csv-upload")
		goto fail
	}

	// Open user's review DB.
	// TODO import into a new db instead?
	db, err = database.OpenReviewDB(basedir.Review(userID, l1, l2))
	if err != nil {
		log.Println(fmt.Errorf("could not open review database (%v-%v): %w", l1, l2, err))
		message = "Something went wrong. Please try again."
		_ = s.ErrorMessage(message, "csv-upload")
		goto fail
	}
	defer db.Close()

	// TODO connect to course db to filter out reviews that are not in the course
	// database?
	if err := replay.Replay(db, file); err != nil {
		if errors.Is(err, replay.ErrHasExistingReviews) {
			message = "Can't import data, because existing reviews were found. Try resetting your progress first."
			_ = s.ErrorMessage(message, "csv-upload")
			goto fail
		}

		log.Println(err)
		message = "Something went wrong. Please try again."
		_ = s.ErrorMessage(message, "csv-upload")
		goto fail
	}

	success = true
	message = "File uploaded."
	_ = s.SuccessMessage(message, "csv-upload")

fail:
	// Don't redirect to settings page.
	// Client might use this API by using fetch.
	sendJSON(w, map[string]any{
		"message": message,
		"success": success,
	})
}
