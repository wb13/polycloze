// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package vocab_size

import (
	"database/sql"
	"fmt"
	"time"
)

type Metric struct {
	Time  time.Time `json:"time"`
	Value float64   `json:"value"`

	initialized bool
}

// Returns vocab size at various points in the given range.
func VocabSize(db *sql.DB, from, to time.Time, step time.Duration) ([]Metric, error) {
	if step < time.Second {
		panic("only supports up to second precision")
	}

	// Initialize return value.
	var series []Metric
	for current := from; current.Before(to); current = current.Add(step) {
		series = append(series, Metric{
			Time: current,
		})
	}

	// Compute vocab size.
	query := `
		SELECT (t - @from)/@step, last_value(v) OVER win
		FROM vocabulary_size_history
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
		return nil, fmt.Errorf("failed to compute vocabulary size: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var i int
		var value float64
		if err := rows.Scan(&i, &value); err != nil {
			return nil, fmt.Errorf("failed to compute vocabulary size: %v", err)
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
				FROM vocabulary_size_history
				WHERE t <= ?
			)
		`
		err = db.QueryRow(query, from.Unix()).Scan(&series[0].Value)
		if err != nil {
			return nil, fmt.Errorf("failed to compute vocabulary size: %v", err)
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
