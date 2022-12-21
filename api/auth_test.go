// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package api

import (
	"database/sql"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-chi/chi/v5"

	"github.com/lggruspe/polycloze/auth"
	"github.com/lggruspe/polycloze/database"
)

// Creates user DB in memory for testing.
// Caller has to close the database after use.
func testDB() *sql.DB {
	db, err := database.OpenAuthDB(":memory:")
	if err != nil {
		panic(err)
	}
	return db
}

// Creates server for testing.
// Do `server.Client()` to create a client.
func testServer(db *sql.DB) *httptest.Server {
	r := chi.NewRouter()
	r.Use(auth.Middleware(db))
	r.HandleFunc("/register", handleRegister)
	r.HandleFunc("/signin", handleSignIn)
	r.HandleFunc("/signout", handleSignOut)
	return httptest.NewServer(r)
}

func resolve(ts *httptest.Server, path string) string {
	return ts.URL + path
}

func TestRegisterNoCSRFToken(t *testing.T) {
	t.Parallel()

	db := testDB()
	defer db.Close()

	ts := testServer(db)
	tc := ts.Client()

	v := url.Values{}
	v.Set("username", "foo")
	v.Set("password", "bar")
	resp, err := tc.PostForm(resolve(ts, "/register"), v)
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	// We know registration failed, because we're back at this page.
	text := doc.Find("form button").First().Text()
	if text != "Register" {
		t.Fatal("expected form button text to be 'Register':", text)
	}
}
