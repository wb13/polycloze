// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package difficulty

import (
	"database/sql"
)

type Difficulty struct {
	Level     int `json:"level"`
	Correct   int `json:"correct"`
	Incorrect int `json:"incorrect"`
	Max       int `json:"max"`
	Min       int `json:"min"`
}

// Returns min difficulty (frequency class of easiest unseen word).
// `db` should have access to `review` and `word` tables.
func minDifficulty(db *sql.DB) int {
	query := `
		SELECT coalesce(min(frequency_class), 0)
		FROM word
		WHERE word NOT IN (
			SELECT item FROM review
		)
	`
	var difficulty int
	_ = db.QueryRow(query).Scan(&difficulty)
	return difficulty
}

// Returns max difficulty (frequency class of hardest unseen word).
// `db` should have access to `review` and `word` tables.
func maxDifficulty(db *sql.DB) int {
	query := `
		SELECT coalesce(max(frequency_class), 0)
		FROM word
		WHERE word NOT IN (
			SELECT item FROM review
		)
	`
	var difficulty int
	_ = db.QueryRow(query).Scan(&difficulty)
	return difficulty
}

// Gets most recent record in difficulty table.
// Returns default values if there is none.
func GetLatest(db *sql.DB) Difficulty {
	min := minDifficulty(db)
	difficulty := Difficulty{
		Level: min,
		Min:   min,
		Max:   maxDifficulty(db),
	}

	query := `SELECT v, correct, incorrect FROM estimated_level`
	_ = db.QueryRow(query).Scan(
		&difficulty.Level,
		&difficulty.Correct,
		&difficulty.Incorrect,
	)
	return difficulty

}
