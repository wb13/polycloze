// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

// For managing database of users.
package database

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

// Upgrades users database to the latest version.
func upgradeUsersDB(db *sql.DB) error {
	return goose.Up(db, "migrations/users")
}

// NOTE Caller has to Close the db.
func OpenUsersDB(path string) (*sql.DB, error) {
	db, err := Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open user database: %v", err)
	}
	if err := upgradeUsersDB(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to upgrade user database: %v", err)
	}

	db.SetMaxOpenConns(1)
	_, _ = db.Exec("PRAGMA journal_mode=WAL")
	return db, nil
}
