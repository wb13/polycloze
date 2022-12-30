// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

// Settings page.
package api

import (
	"log"
	"net/http"

	"github.com/lggruspe/polycloze/auth"
	"github.com/lggruspe/polycloze/sessions"
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
			_ = s.ErrorMessage("Something went wrong. Please try again.")
			goto fail
		}

		id, err := auth.Authenticate(db, username, currentPassword)
		if err != nil {
			log.Println(err)
			_ = s.ErrorMessage("Incorrect password.")
			goto fail
		}

		if err := auth.ChangePassword(db, id, newPassword); err != nil {
			_ = s.ErrorMessage("Something went wrong. Please try again.")
			goto fail
		}

		_ = s.SuccessMessage("Password updated.")
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

	messages, _ := s.Messages()
	s.Data["course"] = course
	s.Data["csrfToken"] = sessions.CSRFToken(s.ID)
	s.Data["messages"] = messages
	renderTemplate(w, "settings.html", s.Data)
}
