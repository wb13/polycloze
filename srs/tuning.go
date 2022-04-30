// Auto-tuning stuff (of spacing algorithm parameters).
package srs

import (
	"database/sql"
	"math"
)

// Computes rate of items in given level that advance to the next level.
// Assumes the specified level and the next level are defined in the database.
func advancementRate(tx *sql.Tx, level int) (float64, error) {
	// The result also includes reviews that use old coefficient values.
	// There wouldn't be enough data if those were also excluded.

	query := `
SELECT count(streak) FROM Review WHERE streak >= ? AND streak <= ?
GROUP BY streak ORDER BY streak ASC
`
	rows, err := tx.Query(query, level, level+1)
	if err != nil {
		return math.NaN(), err
	}
	defer rows.Close()

	var counts []float64
	for rows.Next() {
		var count float64
		if err := rows.Scan(&count); err != nil {
			return count, err
		}
		counts = append(counts, count)
	}

	if len(counts) != 2 {
		panic("expected both levels to be defined")
	}
	if counts[0]*counts[1] == 0.0 {
		panic("expected both levels to be non-empty")
	}

	rate := counts[1] / counts[0]
	if rate < 0.0 || rate > 1.0 {
		panic("expected rate to be between 0 and 1")
	}
	return rate, nil
}

// Gets maximum streak in the database.
// Returns -1 if there is none.
func maxStreak(tx *sql.Tx) int {
	query := `SELECT max(streak) FROM Review`
	row := tx.QueryRow(query)
	streak := -1
	row.Scan(&streak)
	return streak
}

// Updates coefficient.
func updateCoefficient(tx *sql.Tx, level int, coefficient float64) error {
	query := `INSERT INTO Coefficient (streak, coefficient) VALUES (?, ?)`
	_, err := tx.Exec(query, level, coefficient)
	return err
}

// Auto-tunes update coefficients.
func autoTune(tx *sql.Tx) error {
	for i := 0; i < maxStreak(tx)+1; i++ {
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
