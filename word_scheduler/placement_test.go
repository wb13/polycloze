// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package word_scheduler

import (
	"testing"

	"github.com/lggruspe/polycloze/database"
	"github.com/lggruspe/polycloze/utils"
)

func getStats[T database.Querier](q T, frequencyClass int) (int, int) {
	var correct, incorrect int
	query := `SELECT correct, incorrect FROM new_word_stat WHERE frequency_class = ?`
	_ = q.QueryRow(query, frequencyClass).Scan(&correct, &incorrect)
	return correct, incorrect
}

func TestUpdateNewWordStat(t *testing.T) {
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	if correct, incorrect := getStats(db, 1); correct != 0 || incorrect != 0 {
		t.Fatal("expected stats to be empty:", correct, incorrect)
	}

	// Insert stat.
	if err := updateNewWordStat(db, 1, true); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	if correct, incorrect := getStats(db, 1); correct != 1 || incorrect != 0 {
		t.Fatal("expected one correct answer at frequency class 1:", correct, incorrect)
	}

	// Update stats.
	if err := updateNewWordStat(db, 1, false); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	if correct, incorrect := getStats(db, 1); correct != 1 || incorrect != 1 {
		t.Fatal("expected one correct answer and one incorrect answer at frequency class 1:", correct, incorrect)
	}
}

func TestEasiestUnseenEmptyTable(t *testing.T) {
	// Should return 0 instead of null.
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	if class := easiestUnseen(db); class != 0 {
		t.Fatal("expected frequency class of easiest unseen word to be 0 in empty DB:", class)
	}
}

func TestPlacementDefault(t *testing.T) {
	// Should be zero.
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	level := Placement(db)
	if level != 0 {
		t.Fatal("expected initial level to be 0:", level)
	}
}

func ace[T database.Querier](q T, frequencyClass int) {
	// NOTE necessary number of correct answers is probably less than 10
	for i := 0; i < 10; i++ {
		if err := updateNewWordStat(q, frequencyClass, true); err != nil {
			panic(err)
		}
	}
}

func flunk[T database.Querier](q T, frequencyClass int) {
	for i := 0; i < 10; i++ {
		if err := updateNewWordStat(q, frequencyClass, false); err != nil {
			panic(err)
		}
	}
}

func TestPlacementTooEasy(t *testing.T) {
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	ace(db, 1)

	// Should place at next level.
	level := Placement(db)
	if level <= 0 {
		t.Fatal("expected to be placed at higher level:", level)
	}

	// Even if the level has no entry in the DB.
	if correct, incorrect := getStats(db, level); correct != 0 || incorrect != 0 {
		t.Fatal("expected estimated level to have no entry:", correct, incorrect)
	}
}

func TestPlacementTooHard(t *testing.T) {
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	flunk(db, 1)

	// Should stay at current level.
	level := Placement(db)
	if level > 1 {
		t.Fatal("expected to stay at current level (1):", level)
	}
}

func TestPlacementPastCompletedFrequencyClasses(t *testing.T) {
	// Placement test should ignore performance from completed frequency classes.
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	// Insert unseen level 3 word.
	query := `INSERT INTO word (word, frequency_class) VALUES ('unseen', 3)`
	if _, err := db.Exec(query); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	// Even if lower levels are failing but all words at these levels have already
	// been seen...
	flunk(db, 1)
	flunk(db, 2)

	// Placement test should ignore performance at lower levels.
	level := Placement(db)
	if level != 3 {
		t.Fatal("expected to skip to level 3:", level)
	}
}

func TestPlacementGreaterThanOrEqualToEasiestUnseen(t *testing.T) {
	// Placement test level should be >= easiest unseen frequency class.
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	// Insert unseen level 3 word.
	query := `INSERT INTO word (word, frequency_class) VALUES ('unseen', 3)`
	if _, err := db.Exec(query); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	// Scenario: student completed all previous levels without failing or acing them.
	if err := updateNewWordStat(db, 0, true); err != nil {
		t.Fatal("expected err to be nil:", err)
	}
	if err := updateNewWordStat(db, 1, true); err != nil {
		t.Fatal("expected err to be nil:", err)
	}
	if err := updateNewWordStat(db, 2, true); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	// Estimated level >= easiest unseen frequency class.
	level := Placement(db)
	if level != 3 {
		t.Fatal("expected level to be >= 3:", level)
	}
}
