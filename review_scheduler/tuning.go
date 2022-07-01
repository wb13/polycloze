// Auto-tuning stuff (of spacing algorithm parameters).
package review_scheduler

import (
	"database/sql"
	"errors"
	"time"
)

// Initializes interval table.
func initTuner(tx *sql.Tx) error {
	query := `insert or ignore into interval (interval) values (0), (1)`
	_, err := tx.Exec(query)
	return err
}

// Calculates rate of reviews that reach the next level.
// May return 0.925 to avoid setting off auto-tune when there's not enough data.
func advancementRate(correct, incorrect int) float64 {
	if correct+incorrect <= 0 {
		return 0.925
	}
	return float64(correct) / float64(correct+incorrect)
}

// Auto-tunes intervals.
func autoTune(tx *sql.Tx) error {
	if err := initTuner(tx); err != nil {
		return err
	}

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

		if interval <= 1 {
			// Don't change intervals = 0 and 1.
			continue
		}

		rate := advancementRate(correct, incorrect)
		// Target rate is between 90 (to reduce congestion) and 95% (could be higher,
		// but it would be hard to tell if the spacing between levels is too short).
		if rate < 0.90 {
			if err := decreaseInterval(tx, interval); err != nil {
				return err
			}
		} else if rate > 0.95 {
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
		if errors.Is(err, sql.ErrNoRows) {
			panic("unreachable")
		}
		return 0, err
	}
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

	query = `update review set interval = ? where interval = ?`
	_, err := tx.Exec(query, replacement, interval)
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

	query := `insert or ignore into interval (interval) values ?`
	if _, err := tx.Exec(query, 2*max); err != nil {
		return 0, err
	}
	return 2 * max, nil
}

// Returns smallest interval bigger than the specified value.
func nextInterval(tx *sql.Tx, interval time.Duration) (time.Duration, error) {
	query := `select min(interval) from interval where interval > ?`
	row := tx.QueryRow(query, interval)

	var next time.Duration
	if err := row.Scan(&next); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return insertNextInterval(tx, interval)
		}
		return 0, err
	}
	return next, nil
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
