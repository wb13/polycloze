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
	"github.com/lggruspe/polycloze/sessions"
)

type Word struct {
	Word     string    `json:"word"`
	Learned  time.Time `json:"learned"`
	Reviewed time.Time `json:"reviewed"`
	Due      time.Time `json:"due"`
	Strength int       `json:"strength"`
}

type Vocabulary struct {
	Words []Word `json:"words"`
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

	bytes, err := json.Marshal(Vocabulary{Words: results})
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

// Returns map from interval (as number of hours) to strength.
// Use this to compute interval strength.
// The result is not the same as `interval.ROWID`, because there can be gaps in
// rowids.
func queryIntervalStrengths(db *sql.DB) (map[int]int, error) {
	query := `SELECT interval FROM interval ORDER BY interval ASC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query intervals: %v", err)
	}
	defer rows.Close()

	var strength int
	intervals := make(map[int]int)
	for rows.Next() {
		var interval int
		if err := rows.Scan(&interval); err != nil {
			return nil, fmt.Errorf("failed to query intervals: %v", err)
		}
		intervals[interval] = strength
		strength++
	}
	return intervals, nil
}

// Lists words returned by query.
// - limit should be between 10 and 100.
// 	 Silently changes limit if not.
func searchVocabulary(db *sql.DB, limit int, after string, sortBy string) ([]Word, error) {
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

	intervals, err := queryIntervalStrengths(db)
	if err != nil {
		return nil, fmt.Errorf("vocabulary search failed: %v", err)
	}

	query := fmt.Sprintf(`
		SELECT item AS word, learned, reviewed, due, interval
		FROM review JOIN interval USING (interval)
		WHERE item > ?
		ORDER BY %s
		LIMIT ?
	`, sortBy)

	words := make([]Word, 0)
	rows, err := db.Query(query, after, limit)
	if err != nil {
		return nil, fmt.Errorf("vocabulary search failed: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var vocab Word
		var learned, reviewed, due int64
		var interval int
		if err := rows.Scan(&vocab.Word, &learned, &reviewed, &due, &interval); err != nil {
			return nil, fmt.Errorf("vocabulary search failed: %v", err)
		}
		vocab.Learned = time.Unix(learned, 0)
		vocab.Reviewed = time.Unix(reviewed, 0)
		vocab.Due = time.Unix(due, 0)
		vocab.Strength = intervals[interval]
		words = append(words, vocab)
	}
	return words, nil
}
