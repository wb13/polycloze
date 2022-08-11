// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package database

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestUpgrade(t *testing.T) {
	t.Parallel()

	db, _ := sql.Open("sqlite3", ":memory:")
	if err := Upgrade(db); err != nil {
		t.Fatal("expected err to be nil", err)
	}
}

func TestUpgradeTwice(t *testing.T) {
	// Migration should go smoothly both times, even if there are no changes.
	t.Parallel()

	db, _ := sql.Open("sqlite3", ":memory:")
	if err := Upgrade(db); err != nil {
		t.Fatal("expected err to be nil on first upgrade", err)
	}
	if err := Upgrade(db); err != nil {
		t.Fatal("expected err to be nil on second upgrade", err)
	}
}
