package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/lggruspe/polycloze/basedir"
	"github.com/lggruspe/polycloze/buffer"
	"github.com/lggruspe/polycloze/database"
	"github.com/lggruspe/polycloze/flashcards"
	"github.com/lggruspe/polycloze/review_scheduler"
)

type Items struct {
	Items []flashcards.Item `json:"items"`
}

func generateFlashcards(buf *buffer.ItemBuffer, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bytes, err := json.Marshal(Items{Items: buf.TakeMany()})
	if err != nil {
		log.Fatal(err)
	}
	w.Write(bytes)
}

func success() []byte {
	return []byte("{\"success\": true}")
}

func handleReviewUpdate(ig *flashcards.ItemGenerator, w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		panic("expected json body in POST request")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic("something went wrong")
	}

	var reviews Reviews
	if err := json.Unmarshal(body, &reviews); err != nil {
		panic("parsing error")
	}

	session, err := ig.Session()
	if err != nil {
		panic("something went wrong")
	}
	defer session.Close()

	for _, review := range reviews.Reviews {
		review_scheduler.UpdateReview(session, review.Word, review.Correct)
	}
	w.Write(success())
}

// Middleware
func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type")
		next.ServeHTTP(w, r)
	})
}

// NOTE Assumes globalSession has been set using changeLanguages.
func createHandler(w http.ResponseWriter, r *http.Request) {
	buf := globalSession.ItemBuffer
	switch r.Method {
	case "POST":
		handleReviewUpdate(&buf.ItemGenerator, w, r)
	case "GET":
		generateFlashcards(buf, w, r)
	}
}

func handleTest(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		l1 := chi.URLParam(r, "l1")
		l2 := chi.URLParam(r, "l2")
		if err := changeLanguages(db, l1, l2); err != nil {
			log.Fatal(err)
		}
		w.Write([]byte("changed languages"))
	}
}

func Router(config Config) (chi.Router, error) {
	r := chi.NewRouter()

	reviewDb := path.Join(basedir.StateDir, "user", fmt.Sprintf("%v.db", config.L2))
	db, err := database.New(reviewDb)
	if err != nil {
		return r, err
	}
	if err := changeLanguages(db, config.L1, config.L2); err != nil {
		return r, err
	}

	if config.AllowCORS {
		r.Use(cors)
	}
	r.Use(middleware.Logger)
	r.HandleFunc("/", createHandler)
	r.HandleFunc("/options", languageOptions)
	r.HandleFunc("/{l1}/{l2}", handleTest(db))
	return r, nil
}
