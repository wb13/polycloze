// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

// Historical records of user's estimated level.
package history

import (
	"database/sql"
	"fmt"
	"time"
)

// Returns estimated level (frequency class at various points in the given
// range).
func EstimatedLevel(db *sql.DB, from, to time.Time, step time.Duration) ([]Metric[int], error) {
	series := Zeros[int](from, to, step)
	query := `
		SELECT (t - @from)/@step, last_value(v) OVER win
		FROM estimated_level_history
		WHERE t >= @from AND t < @to
		WINDOW win AS (
			PARTITION BY (t - @from)/@step
			ORDER BY id ASC
		)
	`
	rows, err := db.Query(
		query,
		sql.Named("from", from.Unix()),
		sql.Named("to", to.Unix()),
		sql.Named("step", step/time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to compute estimated level: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var i, value int
		if err := rows.Scan(&i, &value); err != nil {
			return nil, fmt.Errorf("failed to compute estimated level: %w", err)
		}
		series[i].Value = value
		series[i].initialized = true
	}

	// Insert missing values.
	if len(series) > 0 && !series[0].initialized {
		query = `
			SELECT coalesce(v, 0)
			FROM (
				SELECT max(id), v
				FROM estimated_level_history
				WHERE t <= ?
			)
		`
		err := db.QueryRow(query, from.Unix()).Scan(&series[0].Value)
		if err != nil {
			return nil, fmt.Errorf("failed to compute estimated level: %w", err)
		}
		series[0].initialized = true
	}

	for i := 1; i < len(series); i++ {
		if !series[i].initialized {
			series[i].initialized = true
			series[i].Value = series[i-1].Value
		}
	}
	return series, nil
}
