// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package api

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/lggruspe/polycloze/auth"
	"github.com/lggruspe/polycloze/basedir"
	"github.com/lggruspe/polycloze/database"
	"github.com/lggruspe/polycloze/difficulty"
	"github.com/lggruspe/polycloze/flashcards"
	"github.com/lggruspe/polycloze/sessions"
	"github.com/lggruspe/polycloze/text"
	"github.com/lggruspe/polycloze/word_scheduler"
)

// Returns predicate to pass to item generator.
func excludeWords(r *http.Request) func(string) bool {
	exclude := make(map[string]bool)
	for _, word := range r.URL.Query()["x"] {
		exclude[text.Casefold(word)] = true
	}
	return func(word string) bool {
		_, found := exclude[text.Casefold(word)]
		return !found
	}
}

// Saves review results to the db.
// Returns an error if it fails to save one or more of the review results.
// The caller may choose to ignore the error.
func saveReviewResults[T database.Querier](q T, reviews []ReviewResult) error {
	var err error
	for _, review := range reviews {
		_err := word_scheduler.UpdateWord(q, review.Word, review.Correct)
		if _err != nil {
			err = _err
		}
	}

	if err != nil {
		return fmt.Errorf("failed to save some reviews: %v", err)
	}
	return err
}

func handleReviewUpdate(con *database.Connection, w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "expected json body in POST request", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not read request.", http.StatusInternalServerError)
		return
	}

	var reviews FlashcardsRequest
	if err := parseJSON(w, body, &reviews); err != nil {
		return
	}

	if len(reviews.Reviews) > 0 {
		// Check csrf token in HTTP headers.
		if !sessions.CheckCSRFToken(s.ID, r.Header.Get("X-CSRF-Token")) {
			http.Error(w, "Forbidden.", http.StatusForbidden)
			return
		}

		if err := saveReviewResults(con, reviews.Reviews); err != nil {
			log.Println(err)
		}
	}

	sendJSON(w, FlashcardsResponse{
		Difficulty: difficulty.GetLatest(con),
	})
}

func handleFlashcards(w http.ResponseWriter, r *http.Request) {
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
		log.Println(fmt.Errorf("could not open review database (%v-%v): %v", l1, l2, err))
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Create database connection with access to review and course DB.
	hook := database.AttachCourse(basedir.Course(l1, l2))
	con, err := database.NewConnection(db, r.Context(), hook)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	defer con.Close()

	switch r.Method {
	case "POST":
		handleReviewUpdate(con, w, r, s)
	case "GET":
		items := flashcards.Get(con, getN(r), excludeWords(r))
		sendJSON(w, FlashcardsResponse{
			Items:      items,
			Difficulty: difficulty.GetLatest(con),
		})
	}
}
