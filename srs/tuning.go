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
	query := `SELECT item, level FROM Review ORDER BY id ASC`
	rows, err := tx.Query(query)
	if err != nil {
		return math.NaN(), err
	}
	defer rows.Close()

	increased := 0
	decreased := 0
	// NOTE items that stayed at the same level after review (because the student
	// crammed) are not counted

	fromLevel := make(map[string]bool)
	for rows.Next() {
		var item string
		var lv int
		if err := rows.Scan(&item, &lv); err != nil {
			return math.NaN(), err
		}

		if lv < level && fromLevel[item] {
			decreased++
			fromLevel[item] = false
		} else if lv == level {
			fromLevel[item] = true
		} else if lv > level+1 {
			increased++
			fromLevel[item] = false
		}
	}

	if decreased == 0 && increased == 0 {
		return 0.925, nil
	}
	return float64(increased) / float64(decreased+increased), nil
}

// Returns maximum review level.
func maxLevel(tx *sql.Tx) (int, error) {
	query := `SELECT max(level) FROM Review`
	row := tx.QueryRow(query)
	var level int
	if err := row.Scan(&level); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	return level, nil
}

// Gets level coefficient.
// Returns the default value (2.0) on error.
func getCoefficient(tx *sql.Tx, level int) float64 {
	query := `SELECT coefficient FROM UpdatedCoefficient WHERE level = ?`
	row := tx.QueryRow(query, level)
	coefficient := 2.0
	row.Scan(&coefficient)
	return coefficient
}

// Sets new coefficient for level.
func setCoefficient(tx *sql.Tx, level int, coefficient float64) error {
	query := `INSERT INTO Coefficient (level, coefficient) VALUES (?, ?)`
	_, err := tx.Exec(query, level, coefficient)
	return err
}

// Auto-tunes update coefficients.
func autoTune(tx *sql.Tx) error {
	max, err := maxLevel(tx)
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
			err = setCoefficient(tx, i, (1+coefficient)/2)
		} else if rate > 0.95 {
			err = setCoefficient(tx, i, coefficient*2)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
