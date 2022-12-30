// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

// For managing user DBs.
package database

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

// Upgrades user DB to the latest version.
func upgradeUserDB(db *sql.DB) error {
	if err := goose.Up(db, "migrations/users"); err != nil {
		return fmt.Errorf("failed to upgrade user database: %w", err)
	}
	return nil
}

// Opens database for one user.
// The caller has to Close the db.
func OpenUserDB(path string) (*sql.DB, error) {
	db, err := Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open user database: %w", err)
	}
	if err := upgradeUserDB(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to open user database: %w", err)
	}
	return db, nil
}
