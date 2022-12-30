// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package difficulty

import (
	"fmt"

	"github.com/lggruspe/polycloze/database"
)

type Difficulty struct {
	Level     int `json:"level"`
	Correct   int `json:"correct"`
	Incorrect int `json:"incorrect"`
	Max       int `json:"max"`
	Min       int `json:"min"`
}

// Returns min difficulty (frequency class of easiest unseen word).
// `Querier` should have access to `review` and `word` tables.
func minDifficulty[T database.Querier](q T) int {
	query := `
		SELECT coalesce(min(frequency_class), 0)
		FROM word
		WHERE word NOT IN (
			SELECT item FROM review
		)
	`
	var difficulty int
	_ = q.QueryRow(query).Scan(&difficulty)
	return difficulty
}

// Returns max difficulty (frequency class of hardest unseen word).
// `Querier` should have access to `review` and `word` tables.
func maxDifficulty[T database.Querier](q T) int {
	query := `
		SELECT coalesce(max(frequency_class), 0)
		FROM word
		WHERE word NOT IN (
			SELECT item FROM review
		)
	`
	var difficulty int
	_ = q.QueryRow(query).Scan(&difficulty)
	return difficulty
}

// Gets most recent record in difficulty table.
// Returns default values if there is none.
func GetLatest[T database.Querier](q T) Difficulty {
	// TODO find way to avoid querying min and max difficulties, because they get
	// slower and slower (eventually O(nlg(n)) even with indexed columns).
	min := minDifficulty(q)
	difficulty := Difficulty{
		Level: min,
		Min:   min,
		Max:   maxDifficulty(q),
	}

	query := `SELECT v, correct, incorrect FROM estimated_level`
	_ = q.QueryRow(query).Scan(
		&difficulty.Level,
		&difficulty.Correct,
		&difficulty.Incorrect,
	)
	return difficulty
}

// Updates difficulty table.
func Update[T database.Querier](q T, difficulty Difficulty) error {
	query := `
		INSERT OR REPLACE INTO estimated_level (v, correct, incorrect)
		VALUES (?, ?, ?)
	`
	_, err := q.Exec(
		query,
		difficulty.Level,
		difficulty.Correct,
		difficulty.Incorrect,
	)
	if err != nil {
		return fmt.Errorf("failed to update difficulty table: %w", err)
	}
	return nil
}
