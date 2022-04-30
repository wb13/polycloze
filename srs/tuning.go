// Auto-tuning stuff (of spacing algorithm parameters).
package srs

import (
	"math"
)

// Computes rate of items in given level that advance to the next level.
// Assumes the specified level and the next level are defined in the database.
func advancementRate(ws *WordScheduler, level int) (float64, error) {
	// The result also includes reviews that use old coefficient values.
	// There wouldn't be enough data if those were also excluded.

	query := `
SELECT streak, count(streak) FROM Review WHERE streak >= ? AND streak <= ?
GROUP BY streak ORDER BY streak
`
	rows, err := ws.db.Query(query, level, level+1)
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
	return counts[1] / counts[0], nil
}
