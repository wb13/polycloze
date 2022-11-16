// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/lggruspe/polycloze/activity"
	"github.com/lggruspe/polycloze/auth"
	"github.com/lggruspe/polycloze/basedir"
	"github.com/lggruspe/polycloze/database"
	"github.com/lggruspe/polycloze/sessions"
)

func handleActivity(w http.ResponseWriter, r *http.Request) {
	db := auth.GetDB(r)
	s, err := sessions.ResumeSession(db, w, r)
	if err != nil || !isSignedIn(s) {
		http.NotFound(w, r)
		return
	}

	l1 := chi.URLParam(r, "l1")
	l2 := chi.URLParam(r, "l2")
	if !courseExists(l1, l2) {
		http.NotFound(w, r)
		return
	}

	userID := s.Data["userID"].(int)
	db, err = database.New(basedir.Review(userID, l1, l2))
	if err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	history, err := activity.ActivityHistory(db, time.Now())
	if err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	// Encode JSON
	data := make(map[string][]activity.Activity)
	data["activities"] = history
	bytes, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	// Send JSON
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(bytes); err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
	}
}
