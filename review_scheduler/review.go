// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

// Defines the Review struct for storing rows from the Review table.
package review_scheduler

import (
	"database/sql"
	"time"
)

type Review struct {
	Interval time.Duration // Interval between now and due date
	Reviewed time.Time
}

func (r Review) Due() time.Time {
	return r.Reviewed.Add(r.Interval)
}

func (r Review) Correct() bool {
	return r.Interval > 0
}

// Calculates interval for next review.
func calculateInterval(tx *sql.Tx, review *Review, correct bool, now time.Time) (time.Duration, error) {
	if !correct {
		return 0, nil
	}
	//reviewed := now
	if review != nil {
		if now.Before(review.Due()) {
			// Don't increase interval if the user crammed
			println("user crammed :(")
			return review.Interval, nil
		}
		//reviewed = review.Reviewed
	}

	//interval := now.Sub(reviewed) // this is greater than review.Interval
	// ^ causes interval to grow disproportionately when you spend time away from training?
	interval := review.Interval // Just use existing interval
	return nextInterval(tx, interval)
}

// Computes next review schedule.
// If review is nil, creates Review with default values for initial review.
// now should usually be time.Now.UTC().
func nextReview(tx *sql.Tx, review *Review, correct bool, now time.Time) (Review, error) {
	var r Review
	interval, err := calculateInterval(tx, review, correct, now)
	if err != nil {
		return r, err
	}
	r.Reviewed = now
	r.Interval = interval
	return r, nil
}
