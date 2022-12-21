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
	db, err := database.OpenAuthDB(":memory:")
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

func TestChangePassword(t *testing.T) {
	t.Parallel()
	db := openDB()
	defer db.Close()

	if err := Register(db, "foo", "bar"); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	oldID, err := Authenticate(db, "foo", "bar")
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	// Change password
	if err := ChangePassword(db, oldID, "baz"); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	// Authenticating with old password should fail
	if _, err := Authenticate(db, "foo", "bar"); err == nil {
		t.Fatal("expected sign in with old password to fail")
	}

	// Authenticating with new password shouldn't fail
	newID, err := Authenticate(db, "foo", "baz")
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	// user ID shouldn't change after password change.
	if newID != oldID {
		t.Fatal("expected user ID to not change:", oldID, newID)
	}
}

func TestChangePasswordNonExistentUser(t *testing.T) {
	t.Parallel()
	db := openDB()
	defer db.Close()

	if err := ChangePassword(db, 1, "baz"); err != nil {
		t.Fatal("ChangePassword should not return error if user doesn't exist:", err)
	}
}

func TestChangePasswordStorage(t *testing.T) {
	t.Parallel()
	db := openDB()
	defer db.Close()

	// Register user
	if err := Register(db, "foo", "bar"); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	id, err := Authenticate(db, "foo", "bar")
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	// Change password
	password := "password"
	if err := ChangePassword(db, id, password); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	// Check how the password is stored.
	var hash string
	query := `SELECT password FROM user WHERE username = 'foo'`
	if err := db.QueryRow(query).Scan(&hash); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	if strings.Contains(hash, password) {
		t.Fatal("password should not be stored in plaintext")
	}
	if strings.Contains(password, hash) {
		t.Fatal("password should not be stored in plaintext")
	}
}
