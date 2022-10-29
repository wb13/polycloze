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

func newTemplateData(r *http.Request) (templateData, error) {
	var data templateData
	session, err := auth.GetSession(r)
	if err != nil {
		return data, err
	}
	data.Session = &session
	return data, nil
}

func handleSignIn(w http.ResponseWriter, r *http.Request) {
	// TODO redirect if logged in
	data, err := newTemplateData(r)
	if err != nil {
		data.Message = "Authentication failed."
		goto fail
	}

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
			data.Message = "Incorrect username or password."
			goto fail
		}

		session, err := auth.GetSession(r)
		if err != nil {
			data.Message = "Authentication failed."
			goto fail
		}

		session.Data.UserID = userID
		session.Data.Username = username

		if err := session.Save(w); err != nil {
			data.Message = "Authentication failed."
			goto fail
		}
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

fail:
	if err := renderTemplate(w, "signin.html", &data); err != nil {
		log.Println(err)
	}
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	// TODO redirect if logged in
	data, err := newTemplateData(r)
	if err != nil {
		data.Message = "Something went wrong."
		goto fail
	}

	if r.Method == "POST" {
		db, err := database.OpenUsersDB(basedir.Users())
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		username := r.FormValue("username")
		password := r.FormValue("password")
		if err := auth.Register(db, username, password); err != nil {
			data.Message = "This username is unavailable. Try another one."
			goto fail
		}

		http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
		return
	}

fail:
	if err := renderTemplate(w, "register.html", &data); err != nil {
		log.Println(err)
	}
}

func handleSignOut(w http.ResponseWriter, r *http.Request) {
	// TODO what if not signed in?
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}

	session, _ := auth.GetSession(r) // TODO handle error
	_ = session.Delete(w)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
