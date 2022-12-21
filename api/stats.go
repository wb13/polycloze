// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/lggruspe/polycloze/auth"
	"github.com/lggruspe/polycloze/basedir"
	"github.com/lggruspe/polycloze/database"
	"github.com/lggruspe/polycloze/history"
	"github.com/lggruspe/polycloze/sessions"
	"github.com/lggruspe/polycloze/vocab_size"
)

// If upgrade is non-empty, upgrades the database.
func queryInt(path, query string, upgrade ...bool) (int, error) {
	var result int

	db, err := database.Open(path)
	if err != nil {
		return 0, fmt.Errorf("could not open db (%v) for query (%v): %v", path, query, err)
	}
	defer db.Close()

	if len(upgrade) > 0 {
		if err := database.UpgradeReviewDB(db); err != nil {
			return result, err
		}
	}

	row := db.QueryRow(query)
	err = row.Scan(&result)
	return result, err
}

// Total count of words in course.
func CountTotal(l1, l2 string) (int, error) {
	return queryInt(basedir.Course(l1, l2), `select count(*) from word`)
}

func handleStatsActivity(w http.ResponseWriter, r *http.Request) {
	db := auth.GetDB(r)
	s, err := sessions.ResumeSession(db, w, r)
	if err != nil || !s.IsSignedIn() {
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
	db, err = database.OpenReviewDB(basedir.Review(userID, l1, l2))
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	result, err := history.Summarize(
		db,
		getFrom(r),
		getTo(r),
		getStep(r),
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	sendJSON(w, map[string]any{
		"activity": result,
		// TODO use unix timestamps?
	})
}

// Responds with user's vocab size over time.
func handleStatsVocab(w http.ResponseWriter, r *http.Request) {
	db := auth.GetDB(r)
	s, err := sessions.ResumeSession(db, w, r)
	if err != nil || !s.IsSignedIn() {
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
	db, err = database.OpenReviewDB(basedir.Review(userID, l1, l2))
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	result, err := vocab_size.VocabSize(
		db,
		getFrom(r),
		getTo(r),
		getStep(r),
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	sendJSON(w, map[string]any{
		"vocabSize": result,
		// TODO use unix timestamps?
	})
}

// Gets `from` UNIX timestamp from URL search params.
// Default value: last week.
func getFrom(r *http.Request) time.Time {
	q := r.URL.Query()
	v := q.Get("from")

	parsed, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return time.Now().AddDate(0, 0, -7)
	}
	return time.Unix(parsed, 0)
}

// Gets `to` UNIX timestamp from URL search params.
// Default value: now.
func getTo(r *http.Request) time.Time {
	q := r.URL.Query()
	v := q.Get("to")

	parsed, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return time.Now()
	}
	return time.Unix(parsed, 0)
}

// Gets `step` size (number of seconds) from URL search params.
func getStep(r *http.Request) time.Duration {
	q := r.URL.Query()
	v := q.Get("step")

	parsed, err := strconv.ParseInt(v, 10, 64)
	if err != nil || parsed < 1 {
		// Default return value if invalid arg.
		return 24 * time.Hour
	}
	return time.Duration(parsed) * time.Second
}
