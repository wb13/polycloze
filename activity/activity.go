// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package activity

import (
	"database/sql"
	"fmt"
	"time"
)

type Activity struct {
	Forgotten    int `json:"forgotten"`
	Unimproved   int `json:"unimproved"`
	Crammed      int `json:"crammed"`
	Learned      int `json:"learned"`
	Strengthened int `json:"strengthened"`
}

// Returns user activity over the past year.
// The result is a slice of length 366.
// result[i]: Activity i days ago.
func ActivityHistory(db *sql.DB, now time.Time) ([]Activity, error) {
	query := `
		SELECT forgotten, unimproved, crammed, learned, strengthened, ? - days_since_epoch AS i
		FROM activity
		WHERE i >= 0 AND i < 366
	`

	rows, err := db.Query(query, now.Unix()/60/60/24)
	if err != nil {
		return nil, fmt.Errorf("failed to get activity history: %v", err)
	}
	defer rows.Close()

	activities := make([]Activity, 366)
	for rows.Next() {
		var i int
		var activity Activity
		if err := rows.Scan(
			&activity.Forgotten,
			&activity.Unimproved,
			&activity.Crammed,
			&activity.Learned,
			&activity.Strengthened,
			&i,
		); err != nil {
			return nil, fmt.Errorf("failed to get activity history: %v", err)
		}
		activities[i] = activity
	}
	return activities, nil
}
