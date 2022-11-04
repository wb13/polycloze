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
		http.NotFound(w, r)
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
			log.Println(username)
			log.Println(currentPassword)
			log.Println(newPassword)
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
	data["csrfToken"] = sessions.CSRFToken(s.ID)
	renderTemplate(w, "settings.html", data)
}
