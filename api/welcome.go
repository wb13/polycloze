// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"

	"github.com/lggruspe/polycloze/auth"
	"github.com/lggruspe/polycloze/basedir"
	"github.com/lggruspe/polycloze/sessions"
)

// Shows welcome page to new user.
func handleWelcome(w http.ResponseWriter, r *http.Request) {
	// TODO don't show if user isn't new
	db := auth.GetDB(r)
	s, err := sessions.ResumeSession(db, w, r)
	if err != nil || !isSignedIn(s) {
		http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
		return
	}

	// Read and parse courses.json to get list of courses.
	path := filepath.Join(basedir.StateDir, "courses.json")
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	var data map[string][]Course
	if err := json.Unmarshal(bytes, &data); err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	// Extract courses from data.
	courses, ok := data["courses"]
	if !ok {
		log.Println("malformed courses.json")
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	// Get L1 and L2 languages.
	var l1Options []Language
	var l2Options []Language
	l1Visited := make(map[string]bool)
	l2Visited := make(map[string]bool)

	for _, course := range courses {
		if _, ok := l1Visited[course.L1.Code]; !ok {
			l1Options = append(l1Options, course.L1)
			l1Visited[course.L1.Code] = true
		}
		if _, ok := l2Visited[course.L2.Code]; !ok {
			l2Options = append(l2Options, course.L2)
			l2Visited[course.L2.Code] = true
		}
	}

	// Sort languages by code.
	sort.Sort(ByCode(l1Options))
	sort.Sort(ByCode(l2Options))

	// Set template data.
	s.Data["csrfToken"] = sessions.CSRFToken(s.ID)
	s.Data["l1Options"] = l1Options
	s.Data["l2Options"] = l2Options
	s.Data["courses"] = courses
	renderTemplate(w, "welcome.html", s.Data)
}
