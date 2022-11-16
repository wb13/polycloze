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

// Returns how many days ago the oldest review was, if it's not too old.
// Bounds the result between 1 and 366.
func oldestActivity(db *sql.DB, now time.Time) (int, error) {
	today := now.Unix() / 60 / 60 / 24
	query := `SELECT coalesce(max(? - days_since_epoch), 0) FROM activity`
	var age int
	if err := db.QueryRow(query, today).Scan(&age); err != nil {
		return 0, fmt.Errorf("failed to get oldest review: %v", err)
	}

	if age <= 0 {
		age = 1
	} else if age > cutoff {
		age = cutoff
	}
	return age, nil
}

// Returns user activity over the past year.
// result[i]: Activity i days ago.
func ActivityHistory(db *sql.DB, now time.Time) ([]Activity, error) {
	query := `
		SELECT forgotten, unimproved, crammed, learned, strengthened, ? - days_since_epoch AS i
		FROM activity
		WHERE i >= 0 AND i <= ?
	`

	rows, err := db.Query(query, now.Unix()/60/60/24, cutoff)
	if err != nil {
		return nil, fmt.Errorf("failed to get activity history: %v", err)
	}
	defer rows.Close()

	var activities []Activity

	if age, err := oldestActivity(db, now); err != nil {
		return nil, fmt.Errorf("failed to get activity history: %v", err)
	} else {
		activities = make([]Activity, age+1)
	}
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
