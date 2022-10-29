// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

// auth-related handlers.
package api

import (
	"log"
	"net/http"

	"github.com/lggruspe/polycloze/auth"
	"github.com/lggruspe/polycloze/basedir"
	"github.com/lggruspe/polycloze/database"
)

func handleSignIn(w http.ResponseWriter, r *http.Request) {
	var message string
	if r.Method == "POST" {
		db, err := database.OpenUsersDB(basedir.Users())
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		username := r.FormValue("username")
		password := r.FormValue("password")

		userID, err := auth.Authenticate(db, username, password)
		if err != nil {
			message = "Incorrect username or password."
			goto fail
		}

		session := auth.GetSession(r)
		session.Data.UserID = userID
		session.Data.Username = username

		if err := session.Save(w); err != nil {
			message = "Authentication failed."
			goto fail
		}
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

fail:
	if err := templates.ExecuteTemplate(w, "signin.html", message); err != nil {
		log.Println(err)
	}
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	var message string

	if r.Method == "POST" {
		db, err := database.OpenUsersDB(basedir.Users())
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		username := r.FormValue("username")
		password := r.FormValue("password")
		if err := auth.Register(db, username, password); err != nil {
			message = "This username is unavailable. Try another one."
			goto fail
		}

		http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
		return
	}

fail:
	if err := templates.ExecuteTemplate(w, "register.html", message); err != nil {
		log.Println(err)
	}
}
