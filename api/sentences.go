// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package api

import (
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/polycloze/polycloze/basedir"
	"github.com/polycloze/polycloze/database"
	"github.com/polycloze/polycloze/sentences"
)

// Gets limit from URL query.
// Returns a default value if the parameter is missing, invalid or too big.
func getSentencesLimit(q url.Values) int {
	v := q.Get("limit")
	limit, err := strconv.Atoi(v)
	if err != nil {
		return 10
	}
	if limit <= 0 {
		return 1
	}
	if limit > 1000 {
		return 1000
	}
	return limit
}

// Returns random sentence from course.
func handleSentences(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	l1 := q.Get("l1")
	l2 := q.Get("l2")
	if l1 == "" || l2 == "" {
		http.Error(w, "invalid course languages", http.StatusBadRequest)
		return
	}

	db, err := database.Open(basedir.Course(l1, l2))
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	limit := getSentencesLimit(q)
	result, err := sentences.RandomSentences(db, limit)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	sendJSON(w, map[string][]sentences.Sentence{
		"sentences": result,
	})
}
