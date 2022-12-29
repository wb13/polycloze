// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package api

import (
	"bufio"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/lggruspe/polycloze/auth"
	"github.com/lggruspe/polycloze/sessions"
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
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	if header.Header.Get("Content-Type") != "text/csv" {
		http.Error(w, "expected CSV file upload", http.StatusBadRequest)
		return
	}

	if isTooBig(header.Size) {
		// TODO also check in the client
		http.Error(w, "file too big (> 8MB)", http.StatusBadRequest)
		return
	}

	// Read uploaded file.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		log.Println(scanner.Text())
	}

	// Respond.
	_, _ = w.Write([]byte(":)"))
}
