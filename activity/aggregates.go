// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package activity

import (
	"database/sql"
	"fmt"
	"time"
)

const cutoff = 366

// Aggregates counts for old reviews (> 1 year).
func AggregateOld(db *sql.DB, now time.Time) (Activity, error) {
	today := now.Unix() / 60 / 60 / 24
	query := `
		SELECT
			coalesce(sum(forgotten), 0),
			coalesce(sum(unimproved), 0),
			coalesce(sum(crammed), 0),
			coalesce(sum(learned), 0),
			coalesce(sum(strengthened), 0)
		FROM activity
		WHERE ? - days_since_epoch > ?
	`

	var activity Activity
	row := db.QueryRow(query, today, cutoff)
	err := row.Scan(
		&activity.Forgotten,
		&activity.Unimproved,
		&activity.Crammed,
		&activity.Learned,
		&activity.Strengthened,
	)
	if err != nil {
		return activity, fmt.Errorf("failed to aggregate old stats: %v", err)
	}
	return activity, nil
}
