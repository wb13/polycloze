// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package api

import (
	"database/sql"
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
