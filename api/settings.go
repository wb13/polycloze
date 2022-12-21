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
	var data map[string]any
	db := auth.GetDB(r)
	s, err := sessions.ResumeSession(db, w, r)
	if err != nil || !isSignedIn(s) {
		http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
		return
	}
	data = s.Data

	if r.Method == "POST" {
		username := s.Data["username"].(string)
		currentPassword := r.FormValue("current-password")
		newPassword := r.FormValue("new-password")
		csrfToken := r.FormValue("csrf-token")

		if !sessions.CheckCSRFToken(s.ID, csrfToken) {
			data["message"] = "Something went wrong. Please try again."
			goto fail
		}

		id, err := auth.Authenticate(db, username, currentPassword)
		if err != nil {
			log.Println(err)
			data["message"] = "Incorrect password."
			goto fail
		}

		if err := auth.ChangePassword(db, id, newPassword); err != nil {
			data["message"] = "Something went wrong. Please try again."
			goto fail
		}

		data["message"] = "Password updated."
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

	data["course"] = course
	data["csrfToken"] = sessions.CSRFToken(s.ID)
	renderTemplate(w, "settings.html", data)
}
