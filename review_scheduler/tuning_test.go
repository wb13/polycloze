// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package review_scheduler

import (
	"testing"
	"time"

	"github.com/lggruspe/polycloze/utils"
)

func TestInsertInterval(t *testing.T) {
	// Intervals should be stored as number of hours.
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	if err := insertInterval(tx, time.Hour); err != nil {
		t.Fatal("expected err to be nil:", err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	var interval int64
	query := `SELECT max(interval) FROM interval`
	if err := db.QueryRow(query).Scan(&interval); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	if interval != 1 {
		t.Fatal("expected `interval` to be equal to 1:", interval)
	}
}

func TestMaxInterval(t *testing.T) {
	// Should return 0 if there are no intervals in the database.
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	max, err := maxInterval(tx)
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	if max != 0 {
		t.Fatal("expected max interval to be 0:", err)
	}
}
