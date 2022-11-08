// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

// Auto-tuning stuff (of spacing algorithm parameters).
package review_scheduler

import (
	"database/sql"
	"time"

	"github.com/lggruspe/polycloze/wilson"
)

const day time.Duration = 24 * time.Hour

// Auto-tunes intervals.
func autoTune(tx *sql.Tx) error {
	query := `SELECT interval, correct, incorrect FROM interval ORDER BY interval ASC`
	rows, err := tx.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var interval time.Duration
		var correct, incorrect int
		if err := rows.Scan(&interval, &correct, &incorrect); err != nil {
			return err
		}
		interval *= time.Hour

		if interval <= day {
			// Don't change intervals = 0 and 1 day.
			continue
		}

		if wilson.IsTooHard(correct, incorrect) {
			if err := decreaseInterval(tx, interval); err != nil {
				return err
			}
		} else if wilson.IsTooEasy(correct, incorrect) {
			if err := increaseInterval(tx, interval); err != nil {
				return err
			}
		}
	}
	return nil
}

// Returns biggest interval smaller than the specified value.
func previousInterval(tx *sql.Tx, interval time.Duration) (time.Duration, error) {
	if interval <= day {
		return 0, nil
	}
	query := `select max(interval) from interval where interval < ?`
	row := tx.QueryRow(query, int64(interval.Hours()))

	var prev time.Duration
	if err := row.Scan(&prev); err != nil {
		return 0, err
	}
	// NOTE Assumes the query never returns null.
	return prev * time.Hour, nil
}

// Check if interval already exists in db.
func alreadyExists(tx *sql.Tx, interval time.Duration) (bool, error) {
	query := `select * from interval where interval = ?`
	rows, err := tx.Query(query, int64(interval.Hours()))
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		return true, nil
	}
	return false, nil
}

// Replaces interval value in `interval` and `review` tables.
// Assumes replacement isn't already in the interval table.
func replaceInterval(tx *sql.Tx, interval, replacement time.Duration) error {
	query := `
		UPDATE interval SET interval = ?, correct = 0, incorrect = 0
		WHERE interval = ?
	`
	_, err := tx.Exec(query, int64(replacement.Hours()), int64(interval.Hours()))
	if err != nil {
		return err
	}

	query = `
		UPDATE review SET
			interval = ?,
			due = due + (? - interval) * 3600
		WHERE interval = ?
	`
	_, err = tx.Exec(
		query,
		int64(replacement.Hours()),
		int64(replacement.Hours()),
		int64(interval.Hours()),
	)
	return err
}

func replaceWithExistingInterval(tx *sql.Tx, interval, replacement time.Duration) error {
	query := `DELETE FROM interval WHERE interval = ?`
	if _, err := tx.Exec(query, int64(interval.Hours())); err != nil {
		return err
	}

	query = `UPDATE review SET interval = ? WHERE interval = ?`
	_, err := tx.Exec(query, int64(replacement.Hours()), int64(interval.Hours()))
	return err
}

func decreaseInterval(tx *sql.Tx, interval time.Duration) error {
	if interval <= day {
		return nil
	}

	prev, err := previousInterval(tx, interval)
	if err != nil {
		return err
	}
	mid := (prev + interval) / 2

	if exists, err := alreadyExists(tx, mid); err != nil {
		return err
	} else if exists {
		return replaceWithExistingInterval(tx, interval, mid)
	}
	return replaceInterval(tx, interval, mid)
}

// Returns the largest interval in the database.
func maxInterval(tx *sql.Tx) (time.Duration, error) {
	var max time.Duration
	query := `select max(interval) from interval`
	err := tx.QueryRow(query).Scan(&max)
	return max * time.Hour, err
}

// Creates record for interval if it doesn't already exist.
func insertInterval(tx *sql.Tx, interval time.Duration) error {
	query := `insert or ignore into interval (interval) values (?)`
	_, err := tx.Exec(query, int64(interval.Hours()))
	return err
}

// Inserts all needed intervals to increase the specified interval.
func insertMissingIntervals(tx *sql.Tx, interval time.Duration) error {
	max, err := maxInterval(tx)
	if err != nil {
		return err
	}

	if max > interval {
		return nil
	}

	next := 2 * max
	if next <= 0 {
		next = day
	}

	for next <= interval {
		if err := insertInterval(tx, next); err != nil {
			return err
		}
		next *= 2
	}
	// Make sure that a larger interval exists
	return insertInterval(tx, next)
}

// Returns smallest interval bigger than the specified value.
func nextInterval(tx *sql.Tx, interval time.Duration) (time.Duration, error) {
	if err := insertMissingIntervals(tx, interval); err != nil {
		return 0, err
	}

	query := `select min(interval) from interval where interval > ?`
	row := tx.QueryRow(query, int64(interval.Hours()))

	var next time.Duration
	err := row.Scan(&next)
	return next * time.Hour, err
}

func increaseInterval(tx *sql.Tx, interval time.Duration) error {
	next, err := nextInterval(tx, interval)
	if err != nil {
		return err
	}
	mid := (interval + next) / 2

	if exists, err := alreadyExists(tx, mid); err != nil {
		return err
	} else if exists {
		return replaceWithExistingInterval(tx, interval, mid)
	}
	return replaceInterval(tx, interval, mid)
}

// Updates interval table.
func updateIntervalStats(tx *sql.Tx, review *Review, correct bool) error {
	var interval time.Duration = 0
	if review != nil {
		interval = review.Interval
	}

	// Update interval
	query := `update interval set correct = correct + 1 where interval = ?`
	if !correct {
		query = `update interval set incorrect = incorrect + 1 where interval = ?`
	}
	_, err := tx.Exec(query, int64(interval.Hours()))
	return err
}
