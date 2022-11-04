// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/lggruspe/polycloze/auth"
	"github.com/lggruspe/polycloze/basedir"
	"github.com/lggruspe/polycloze/database"
	"github.com/lggruspe/polycloze/review_scheduler"
	"github.com/lggruspe/polycloze/sessions"
)

type vocabularyItem struct {
	Word     string    `json:"word"`
	Reviewed time.Time `json:"reviewed"`
	Due      time.Time `json:"due"`
	Strength int
}

type vocabulary struct {
	Results []vocabularyItem `json:"results"`
}

func handleVocabulary(w http.ResponseWriter, r *http.Request) {
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

	q := r.URL.Query()
	results, err := searchVocabulary(db, getLimit(q), getAfter(q), getSortBy(q))
	if err != nil {
		log.Println(fmt.Errorf("search error: %v", err))
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	bytes, err := json.Marshal(vocabulary{Results: results})
	if err != nil {
		log.Fatal(err)
	}
	if _, err := w.Write(bytes); err != nil {
		log.Println(err)
	}
}

// Gets limit from URL query.
// If the limit is not in the URL query or is invalid, returns the default (20).
func getLimit(q url.Values) int {
	v := q.Get("limit")
	if limit, err := strconv.Atoi(v); err == nil {
		return limit
	}
	return 10
}

// Gets 'after' from URL query.
func getAfter(q url.Values) string {
	return q.Get("after")
}

// Checks if `sortBy` value is valid.
func isValidSortBy(sortBy string) bool {
	switch sortBy {
	case "word":
		fallthrough
	case "reviewed":
		fallthrough
	case "due":
		fallthrough
	case "strength":
		return true
	default:
		return false
	}
}

// Gets 'sortBy' from URL query.
// If `sortBy` is not in the URL query or is invalid, returns "word".
func getSortBy(q url.Values) string {
	sortBy := q.Get("sortBy")
	if isValidSortBy(sortBy) {
		return sortBy
	}
	return "word"
}

// Lists words returned by query.
// - limit should be between 10 and 100.
// 	 Silently changes limit if not.
func searchVocabulary(db *sql.DB, limit int, after string, sortBy string) ([]vocabularyItem, error) {
	// Cap limit.
	if limit < 10 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	if !isValidSortBy(sortBy) {
		panic(fmt.Errorf("invalid sortBy value: %v", sortBy))
	}

	query := fmt.Sprintf(`
		SELECT item AS word, reviewed, due, interval.rowid AS strength
		FROM review JOIN interval USING (interval)
		WHERE item > ?
		ORDER BY %s
		LIMIT ?
	`, sortBy)

	var words []vocabularyItem
	rows, err := db.Query(query, after, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var vocab vocabularyItem
		var reviewed, due string
		if err := rows.Scan(&vocab.Word, &reviewed, &due, &vocab.Strength); err != nil {
			return nil, err
		}

		vocab.Reviewed, err = review_scheduler.ParseTimestamp(reviewed)
		if err != nil {
			continue
		}
		vocab.Due, err = review_scheduler.ParseTimestamp(due)
		if err != nil {
			continue
		}
		words = append(words, vocab)
	}
	return words, nil
}
