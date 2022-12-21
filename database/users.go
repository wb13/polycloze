// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

// For managing authentication database.
package database

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

// Upgrades auth database to the latest version.
func upgradeAuthDB(db *sql.DB) error {
	return goose.Up(db, "migrations/auth")
}

// NOTE Caller has to Close the db.
func OpenAuthDB(path string) (*sql.DB, error) {
	db, err := Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open auth database: %v", err)
	}
	if err := upgradeAuthDB(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to upgrade auth database: %v", err)
	}

	db.SetMaxOpenConns(1)
	_, _ = db.Exec("PRAGMA journal_mode=WAL")
	return db, nil
}
