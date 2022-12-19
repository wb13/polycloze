// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package difficulty

import (
	"testing"

	"github.com/lggruspe/polycloze/utils"
)

func TestMinDifficultyEmptyTable(t *testing.T) {
	// Should return 0.
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	if value := minDifficulty(db); value != 0 {
		t.Fatal("expected minimum difficulty to be 0 in empty DB:", value)
	}
}

func TestMaxDifficultyEmptyTable(t *testing.T) {
	// Should return 0.
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	if value := maxDifficulty(db); value != 0 {
		t.Fatal("expected maximum difficulty to be 0 in empty DB:", value)
	}
}

func TestGetLatestEmpty(t *testing.T) {
	// Should return 0 on all stats.
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	output := GetLatest(db)
	expected := Difficulty{}

	if output != expected {
		t.Fatal("expected all fields to be 0 initially:", output)
	}
}

func TestUpdate(t *testing.T) {
	// New value should replace old one.
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	first := Difficulty{
		Level:     1,
		Min:       1,
		Max:       1,
		Correct:   1,
		Incorrect: 1,
	}
	second := Difficulty{
		Level:     2,
		Min:       2,
		Max:       2,
		Correct:   2,
		Incorrect: 2,
	}

	// Save first value.
	if err := Update(db, first); err != nil {
		t.Fatal("expected err to be nil:", err)
	}
	out := GetLatest(db)
	equal := first.Level == out.Level &&
		first.Correct == out.Correct &&
		first.Incorrect == out.Incorrect
	// We ignore Min and Max, because they don't get saved in the DB.

	if !equal {
		t.Fatal("expected latest to be the most recently saved:", out, first)
	}

	// Save second value.
	if err := Update(db, second); err != nil {
		t.Fatal("expected err to be nil:", err)
	}
	out = GetLatest(db)
	equal = second.Level == out.Level &&
		second.Correct == out.Correct &&
		second.Incorrect == out.Incorrect
	// We ignore Min and Max, because they don't get saved in the DB.

	if !equal {
		t.Fatal("expected latest to be the most recently saved:", out, second)
	}
}

func TestGetLatestLevelBelowMinLevel(t *testing.T) {
	// Level should be bounded below by Min.
	// We don't have to test if Level is bounded above by Max, because it
	// doesn't matter if Level > Max (word scheduler would look for easier words
	// in that scenario.
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	// Insert a word.
	query := `INSERT INTO word (word, frequency_class) VALUES ('unseen', 3)`
	if _, err := db.Exec(query); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	// Default level is normally 0, but it should now be 3, because that's the
	// lowest difficulty among unseen words.
	if out := GetLatest(db); out.Level != 3 {
		t.Fatal("expected level >= min (3):", out)
	}
}

func TestGetLatestLevelAboveMinLevel(t *testing.T) {
	// If level > min, GetLatest shouldn't change it.
	t.Parallel()

	db := utils.TestingDatabase()
	defer db.Close()

	// Insert a word.
	query := `INSERT INTO word (word, frequency_class) VALUES ('unseen', 3)`
	if _, err := db.Exec(query); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	if err := Update(db, Difficulty{Level: 5}); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	if out := GetLatest(db); out.Level != 5 {
		t.Fatal("expected level to be 5:", out)
	}
}
