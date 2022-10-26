// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lggruspe/polycloze/basedir"
	"github.com/lggruspe/polycloze/database"
	"github.com/lggruspe/polycloze/flashcards"
	"github.com/lggruspe/polycloze/logger"
	"github.com/lggruspe/polycloze/text"
	"github.com/lggruspe/polycloze/word_scheduler"
)

type Items struct {
	Items []flashcards.Item `json:"items"`
}

func getN(r *http.Request) int {
	n := 10
	q := r.URL.Query()
	if i, err := strconv.Atoi(q.Get("n")); err == nil && i > 0 {
		n = i
	}
	return n
}

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

func generateFlashcards(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	l1 := chi.URLParam(r, "l1")
	l2 := chi.URLParam(r, "l2")
	hook := database.AttachCourse(basedir.Course(l1, l2))
	items := flashcards.Get(db, getN(r), excludeWords(r), hook)
	bytes, err := json.Marshal(Items{Items: items})
	if err != nil {
		log.Fatal(err)
	}
	if _, err := w.Write(bytes); err != nil {
		log.Println(err)
	}
}

// frequencyClass is taken from student.frequency_class
func success(frequencyClass int) []byte {
	return []byte(fmt.Sprintf("{\"success\": true, \"frequencyClass\": %v}", frequencyClass))
}

func handleReviewUpdate(db *sql.DB, l1, l2 string, w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "expected json body in POST request", 400)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal("could not read request body:", err)
	}

	var reviews Reviews
	if err := json.Unmarshal(body, &reviews); err != nil {
		http.Error(w, "could not parse json", 400)
		return
	}

	hook := database.AttachCourse(basedir.Course(l1, l2))
	con, err := database.NewConnection(db, r.Context(), hook)
	if err != nil {
		log.Fatal("could not connect to database:", err)
	}
	defer con.Close()

	var frequencyClass int
	for _, review := range reviews.Reviews {
		err := word_scheduler.UpdateWord(con, review.Word, review.Correct)
		if err != nil {
			log.Printf("failed to update word: '%v'\n\t%v\n", review.Word, err.Error())
		}
		frequencyClass = word_scheduler.PreferredDifficulty(con)
		_ = logger.LogReview(basedir.Log(l1, l2), review.Correct, review.Word)
	}

	if _, err := w.Write(success(frequencyClass)); err != nil {
		log.Println(err)
	}
}

// Middleware
func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type")
		next.ServeHTTP(w, r)
	})
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	l1 := chi.URLParam(r, "l1")
	l2 := chi.URLParam(r, "l2")

	db, err := database.New(basedir.Review(l1, l2))
	if err != nil {
		log.Fatal(fmt.Errorf("could not open review database (%v-%v): %v", l1, l2, err))
	}
	defer db.Close()

	switch r.Method {
	case "POST":
		handleReviewUpdate(db, l1, l2, w, r)
	case "GET":
		generateFlashcards(db, w, r)
	}
}

func Router(config Config) (chi.Router, error) {
	r := chi.NewRouter()
	if config.AllowCORS {
		r.Use(cors)
	}
	r.Use(middleware.Logger)
	r.HandleFunc("/", showPage("home.html"))
	r.HandleFunc("/about", showPage("about.html"))
	r.HandleFunc("/study", showPage("study.html"))

	r.Handle("/dist/*", http.StripPrefix("/dist/", serveDist()))
	r.Handle("/public/*", http.StripPrefix("/public/", servePublic()))

	r.HandleFunc("/courses", courseOptions)

	r.HandleFunc("/{l1}/{l2}", createHandler)
	return r, nil
}
