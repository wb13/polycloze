// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package auth

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/lggruspe/polycloze/database"
)

// NOTE Caller should close DB.
func openDB() *sql.DB {
	db, err := database.OpenUsersDB(":memory:")
	if err != nil {
		panic(err)
	}
	return db
}

func TestAuthenticateUnregistered(t *testing.T) {
	t.Parallel()
	db := openDB()
	defer db.Close()

	if _, err := Authenticate(db, "foo", "bar"); err == nil {
		t.Fatal("authentication should fail if user is not registered")
	}
}

func TestAuthenticateRegistered(t *testing.T) {
	t.Parallel()
	db := openDB()
	defer db.Close()

	if err := Register(db, "foo", "bar"); err != nil {
		t.Fatal("initial registration should succeed:", err)
	}
	if _, err := Authenticate(db, "foo", "bar"); err != nil {
		t.Fatal("authentication should succeed if username and password are correct:", err)
	}
}

func TestAuthenticateIncorrectPassword(t *testing.T) {
	t.Parallel()
	db := openDB()
	defer db.Close()

	if err := Register(db, "foo", "bar"); err != nil {
		t.Fatal("initial registration should succeed:", err)
	}
	if _, err := Authenticate(db, "foo", "baz"); err == nil {
		t.Fatal("authentication should fail if password is incorrect:", err)
	}
}

func TestRegisterPasswordStorage(t *testing.T) {
	t.Parallel()
	db := openDB()
	defer db.Close()

	username := "username"
	password := "password"

	if err := Register(db, username, password); err != nil {
		t.Fatal("initial registration should succeed:", err)
	}

	var hash string
	query := `SELECT password FROM user WHERE username = ?`
	if err := db.QueryRow(query, username).Scan(&hash); err != nil {
		t.Fatal("user database should contain an entry for registered user:", err)
	}

	if strings.Contains(hash, password) {
		t.Fatal("password should not be stored in plaintext")
	}
	if strings.Contains(password, hash) {
		t.Fatal("password should not be stored in plaintext")
	}
}

func TestRegisterTakenUsername(t *testing.T) {
	t.Parallel()
	db := openDB()
	defer db.Close()

	if err := Register(db, "foo", "bar"); err != nil {
		t.Fatal("initial registration should succeed:", err)
	}
	if err := Register(db, "foo", "baz"); err == nil {
		t.Fatal("registration should fail if username is already taken")
	}
}

func TestRegisterEmptyUsername(t *testing.T) {
	t.Parallel()
	db := openDB()
	defer db.Close()

	if err := Register(db, "", "password"); err == nil {
		t.Fatal("registration should fail if username is an empty string")
	}
}

func TestRegisterEmptyPassword(t *testing.T) {
	t.Parallel()
	db := openDB()
	defer db.Close()

	if err := Register(db, "username", ""); err != nil {
		t.Fatal("empty string should be allowed as password:", err)
	}
	if _, err := Authenticate(db, "username", ""); err != nil {
		t.Fatal("empty string should be allowed as password:", err)
	}
}
