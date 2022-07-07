// Auto-tuning stuff (of spacing algorithm parameters).
package review_scheduler

import (
	"database/sql"
	"time"
)

const day time.Duration = 86400000000000 // In nanoseconds

func isTooEasy(correct, incorrect int) bool {
	// Threshold can't be too high or the tuner will be too conservative.
	// Wilson(4, 0, z) ~0.8485 is too low, because Wilson(3, 1) < Wilson(1, 0).
	// This would cause the algorithm to auto-tune as soon as the user makes a
	// mistake, so the sample size is always n <= 4.
	// Wilson(5, 0, z) allows the user to make a mistake without immediate re-tuning.

	z := -0.845                                             // lower bound z-score for one-sided 80% confidence level
	return Wilson(correct, incorrect, z) >= Wilson(5, 0, z) // ~0.875
}

func isTooHard(correct, incorrect int) bool {
	// The threshold is Wilson(1, 0, z), because the interval shouldn't be made
	// easier if the user got all reviews right.

	z := -0.845                                            // lower bound z-score for one-sided 80% confidence level
	return Wilson(correct, incorrect, z) < Wilson(1, 0, z) // ~0.5834
}

// Auto-tunes intervals.
func autoTune(tx *sql.Tx) error {
	query := `select interval, correct, incorrect from interval order by interval asc`
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

		if interval <= day {
			// Don't change intervals = 0 and 1 day.
			continue
		}

		if isTooHard(correct, incorrect) {
			if err := decreaseInterval(tx, interval); err != nil {
				return err
			}
		} else if isTooEasy(correct, incorrect) {
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
	row := tx.QueryRow(query, interval)

	var prev time.Duration
	if err := row.Scan(&prev); err != nil {
		return 0, err
	}
	// NOTE Assumes the query never returns null.
	return prev, nil
}

// Check if interval already exists in db.
func alreadyExists(tx *sql.Tx, interval time.Duration) (bool, error) {
	query := `select * from interval where interval = ?`
	rows, err := tx.Query(query, interval)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		return true, nil
	}
	return false, nil
}

// Assumes replacement isn't already in the interval table.
func replaceInterval(tx *sql.Tx, interval, replacement time.Duration) error {
	query := `
update interval set interval = ?, correct = 0, incorrect = 0
where interval = ?
`
	if _, err := tx.Exec(query, replacement, interval); err != nil {
		return err
	}

	query = `
update review set
	interval = ?,
	due = datetime(unixepoch(due) + (? - interval) / 1e9, 'unixepoch')
where interval = ?
`
	_, err := tx.Exec(query, replacement, replacement, interval)
	return err
}

func replaceWithExistingInterval(tx *sql.Tx, interval, replacement time.Duration) error {
	query := `delete from interval where interval = ?`
	if _, err := tx.Exec(query, interval); err != nil {
		return err
	}

	query = `update review set interval = ? where interval = ?`
	_, err := tx.Exec(query, replacement, interval)
	return err
}

func decreaseInterval(tx *sql.Tx, interval time.Duration) error {
	if interval <= 1 {
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
	query := `select max(interval) from interval`
	row := tx.QueryRow(query)

	var max time.Duration
	if err := row.Scan(&max); err != nil {
		return 0, err
	}
	return max, nil
}

// Inserts double of the largest interval in the database, or twice the specified
// interval, whichever's larger.
func insertNextInterval(tx *sql.Tx, interval time.Duration) (time.Duration, error) {
	max, err := maxInterval(tx)
	if err != nil {
		return 0, err
	}

	if interval > max {
		max = interval
	}

	query := `insert or ignore into interval (interval) values (?)`
	if _, err := tx.Exec(query, 2*max); err != nil {
		return 0, err
	}
	return 2 * max, nil
}

// Returns smallest interval bigger than the specified value.
func nextInterval(tx *sql.Tx, interval time.Duration) (time.Duration, error) {
	query := `select min(interval) from interval where interval > ?`
	row := tx.QueryRow(query, interval)

	var next sql.NullInt64
	if err := row.Scan(&next); err != nil {
		return 0, err
	}
	if !next.Valid {
		return insertNextInterval(tx, interval)
	}
	return time.Duration(next.Int64), nil
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
