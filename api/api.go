// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/lggruspe/polycloze/auth"
	"github.com/lggruspe/polycloze/sessions"
)

// Middleware
func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type")
		next.ServeHTTP(w, r)
	})
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	// Check if user is signed in.
	db := auth.GetDB(r)
	s, err := sessions.StartOrResumeSession(db, w, r)

	if err != nil || !isSignedIn(s) {
		http.Redirect(w, r, "/about", http.StatusTemporaryRedirect)
		return
	}

	// Get active course.
	userID := s.Data["userID"].(int)
	l1Code, l2Code, err := getUserActiveCourse(userID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	s.Data["l1Code"] = l1Code
	s.Data["l2Code"] = l2Code
	renderTemplate(w, "home.html", s.Data)
}

func handleAbout(w http.ResponseWriter, r *http.Request) {
	var data map[string]any
	db := auth.GetDB(r)
	if s, err := sessions.StartOrResumeSession(db, w, r); err == nil {
		data = s.Data

		if isSignedIn(s) {
			// Get active course.
			userID := data["userID"].(int)
			l1Code, l2Code, err := getUserActiveCourse(userID)
			if err != nil {
				log.Println(err)
				http.Error(w, "Something went wrong.", http.StatusInternalServerError)
				return
			}
			data["l1Code"] = l1Code
			data["l2Code"] = l2Code
		}
	}
	renderTemplate(w, "about.html", data)
}

func handleStudy(w http.ResponseWriter, r *http.Request) {
	db := auth.GetDB(r)
	s, err := sessions.ResumeSession(db, w, r)
	if err != nil || !isSignedIn(s) {
		http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
		return
	}

	// Get active course.
	userID := s.Data["userID"].(int)
	l1Code, l2Code, err := getUserActiveCourse(userID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	s.Data["l1Code"] = l1Code
	s.Data["l2Code"] = l2Code
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

	// Get active course.
	userID := s.Data["userID"].(int)
	l1Code, l2Code, err := getUserActiveCourse(userID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	s.Data["l1Code"] = l1Code
	s.Data["l2Code"] = l2Code
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
	r.HandleFunc("/about", handleAbout)
	r.HandleFunc("/welcome", handleWelcome)
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
