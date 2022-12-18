// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package api

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/lggruspe/polycloze/auth"
	"github.com/lggruspe/polycloze/basedir"
	"github.com/lggruspe/polycloze/database"
	"github.com/lggruspe/polycloze/difficulty"
	"github.com/lggruspe/polycloze/flashcards"
	"github.com/lggruspe/polycloze/sessions"
	"github.com/lggruspe/polycloze/text"
)

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
	l1 := chi.URLParam(r, "l1")
	l2 := chi.URLParam(r, "l2")

	hook := database.AttachCourse(basedir.Course(l1, l2))
	con, err := database.NewConnection(db, r.Context(), hook)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	defer con.Close()

	items := flashcards.Get(con, getN(r), excludeWords(r))
	sendJSON(w, FlashcardsResponse{
		Items:      items,
		Difficulty: difficulty.GetLatest(con),
	})
}

func handleReviewUpdate(db *sql.DB, w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "expected json body in POST request", http.StatusBadRequest)
		return
	}

	// Check csrf token in HTTP headers.
	if !sessions.CheckCSRFToken(s.ID, r.Header.Get("X-CSRF-Token")) {
		http.Error(w, "Forbidden.", http.StatusForbidden)
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

	l1 := chi.URLParam(r, "l1")
	l2 := chi.URLParam(r, "l2")
	hook := database.AttachCourse(basedir.Course(l1, l2))
	con, err := database.NewConnection(db, r.Context(), hook)
	if err != nil {
		log.Fatal("could not connect to database:", err)
	}
	defer con.Close()

	if err := saveReviewResults(con, reviews.Reviews); err != nil {
		log.Println(err)
	}
	sendJSON(w, FlashcardsResponse{
		Difficulty: difficulty.GetLatest(con),
	})
}

// Middleware
func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type")
		next.ServeHTTP(w, r)
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

	switch r.Method {
	case "POST":
		handleReviewUpdate(db, w, r, s)
	case "GET":
		generateFlashcards(db, w, r)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	db := auth.GetDB(r)
	s, err := sessions.StartOrResumeSession(db, w, r)

	if err != nil || !isSignedIn(s) {
		http.Redirect(w, r, "/about", http.StatusTemporaryRedirect)
		return
	}
	renderTemplate(w, "home.html", s.Data)
}

func handleStudy(w http.ResponseWriter, r *http.Request) {
	db := auth.GetDB(r)
	s, err := sessions.ResumeSession(db, w, r)
	if err != nil || !isSignedIn(s) {
		http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
		return
	}
	s.Data["csrfToken"] = sessions.CSRFToken(s.ID)
	renderTemplate(w, "study.html", s.Data)
}

func handleVocabularyPage(w http.ResponseWriter, r *http.Request) {
	db := auth.GetDB(r)
	s, err := sessions.ResumeSession(db, w, r)
	if err != nil || !isSignedIn(s) {
		http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
		return
	}
	s.Data["csrfToken"] = sessions.CSRFToken(s.ID)
	renderTemplate(w, "vocab.html", s.Data)
}

// db: user DB for authentication
func Router(config Config, db *sql.DB) (chi.Router, error) {
	r := chi.NewRouter()
	if config.AllowCORS {
		r.Use(cors)
	}
	r.Use(middleware.Logger)
	r.Use(auth.Middleware(db))

	r.HandleFunc("/", handleHome)
	r.HandleFunc("/study", handleStudy)
	r.HandleFunc("/vocab", handleVocabularyPage)
	r.HandleFunc("/about", showPage("about.html"))

	r.HandleFunc("/settings", handleSettings)

	r.HandleFunc("/register", handleRegister)
	r.HandleFunc("/signin", handleSignIn)
	r.HandleFunc("/signout", handleSignOut)

	r.Handle("/dist/*", http.StripPrefix("/dist/", serveDist()))
	r.Handle("/public/*", http.StripPrefix("/public/", servePublic()))
	r.Handle("/share/*", http.StripPrefix("/share/", serveShare()))
	r.Handle("/svg/*", http.StripPrefix("/svg/", serveSVG()))

	// serviceworker has to be at the root.
	r.Handle("/serviceworker.js*", http.StripPrefix("/", serveDist()))

	r.HandleFunc("/{l1}/{l2}/vocab", handleVocabulary)
	r.HandleFunc("/api/sentences", handleSentences)

	r.HandleFunc("/api/flashcards/{l1}/{l2}", handleFlashcards)
	r.HandleFunc("/api/stats/activity/{l1}/{l2}", handleStatsActivity)
	r.HandleFunc("/api/stats/vocab/{l1}/{l2}", handleStatsVocab)

	r.HandleFunc("/api/languages", serveLanguagesJSON())
	r.HandleFunc("/api/courses", serveCoursesJSON())
	return r, nil
}
