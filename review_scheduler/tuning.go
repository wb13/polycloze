// Auto-tuning stuff (of spacing algorithm parameters).
package review_scheduler

import (
	"database/sql"
	"time"
)

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
			// TODO decrease interval
		} else if rate > 0.95 {
			// TODO increase interval
		}
	}
	return nil
}
