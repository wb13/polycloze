// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

// For managing review DBs.
package database

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

// Upgrades review DB to the latest version.
func UpgradeReviewDB(db *sql.DB) error {
	if err := goose.Up(db, "migrations/reviews"); err != nil {
		return fmt.Errorf("failed to upgrade review database: %w", err)
	}
	return nil
}

// Opens review database.
// The caller has to Close the db.
func OpenReviewDB(path string) (*sql.DB, error) {
	db, err := Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open review database: %w", err)
	}
	if err := UpgradeReviewDB(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to open review database: %w", err)
	}
	return db, nil
}
