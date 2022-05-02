// Auto-tuning stuff (of spacing algorithm parameters).
package srs

import (
	"database/sql"
	"errors"
	"math"
)

// Calculates rate of reviews that reach the next level.
// Returns 0.925 when no review has ever left the level, to avoid setting off
// auto-tune when there's not enough data.
func advancementRate(tx *sql.Tx, level int) (float64, error) {
	query := `SELECT word, streak FROM Review ORDER BY id ASC`
	rows, err := tx.Query(query)
	if err != nil {
		return math.NaN(), err
	}
	defer rows.Close()

	failed := 0
	advanced := 0

	activeStreak := make(map[string]bool)
	for rows.Next() {
		var word string
		var streak int
		if err := rows.Scan(&word, &streak); err != nil {
			return math.NaN(), err
		}

		if streak == 0 && activeStreak[word] {
			failed++
			activeStreak[word] = false
		} else if streak == level {
			activeStreak[word] = true
		} else if streak == level+1 {
			advanced++
			activeStreak[word] = false
		}
	}

	if failed == 0 && advanced == 0 {
		return 0.925, nil
	}
	return float64(advanced) / float64(failed+advanced), nil
}

// Returns maximum review streak.
func maxStreak(tx *sql.Tx) (int, error) {
	query := `SELECT max(streak) FROM Review`
	row := tx.QueryRow(query)
	var streak int
	if err := row.Scan(&streak); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	return streak, nil
}

// Updates coefficient.
func updateCoefficient(tx *sql.Tx, level int, coefficient float64) error {
	query := `INSERT INTO Coefficient (streak, coefficient) VALUES (?, ?)`
	_, err := tx.Exec(query, level, coefficient)
	return err
}

// Auto-tunes update coefficients.
func autoTune(tx *sql.Tx) error {
	max, err := maxStreak(tx)
	if err != nil {
		return err
	}

	for i := 0; i < max+1; i++ {
		// Target rate is between 90 (to reduce congestion) and 95% (could be higher,
		// but it would be hard to tell if the spacing between levels is too short).

		coefficient := getCoefficient(tx, i)
		rate, err := advancementRate(tx, i)
		if err != nil {
			return err
		}

		if rate < 0.9 {
			err = updateCoefficient(tx, i, (1+coefficient)/2)
		} else if rate > 0.95 {
			err = updateCoefficient(tx, i, coefficient*2)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
