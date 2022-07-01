// Auto-tuning stuff (of spacing algorithm parameters).
package review_scheduler

import (
	"database/sql"
	"errors"
	"math"
)

// Calculates rate of reviews that reach the next level.
// Returns 0.925 when no review has ever left the level, to avoid setting off
// auto-tune when there's not enough data.
func advancementRate(tx *sql.Tx, level int) (float64, error) {
	query := `SELECT item, level FROM review ORDER BY reviewed ASC`
	rows, err := tx.Query(query)
	if err != nil {
		return math.NaN(), err
	}
	defer rows.Close()

	increased := 0
	decreased := 0
	// Items that stayed at the same level after review (because the student
	// crammed) are not counted

	fromLevel := make(map[string]bool)
	for rows.Next() {
		var item string
		var lv int
		if err := rows.Scan(&item, &lv); err != nil {
			return math.NaN(), err
		}

		if level > 0 && fromLevel[item] {
			if lv < level {
				decreased++
			} else if lv > level {
				increased++
			}
		} else if val, ok := fromLevel[item]; level == 0 && (val || !ok) {
			if lv <= 0 {
				decreased++
			} else if lv > level {
				increased++
			}
		}

		if lv == level {
			fromLevel[item] = true
		} else {
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
	query := `SELECT max(level) FROM review`
	row := tx.QueryRow(query)
	var level int
	if err := row.Scan(&level); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	return level, nil
}

// Auto-tunes update coefficients.
func autoTune(tx *sql.Tx) error {
	max, err := maxLevel(tx)
	if err != nil {
		return err
	}

	// Don't update coefficients for 0, because next interval is always 1 day
	// regardless of advancementRate.
	for i := 1; i < max+1; i++ {
		// Target rate is between 90 (to reduce congestion) and 95% (could be higher,
		// but it would be hard to tell if the spacing between levels is too short).

		_, err := advancementRate(tx, i)
		if err != nil {
			return err
		}
		// TODO auto-tune based on advancement rates
	}
	return nil
}
