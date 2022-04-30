// Defines the Review struct for storing rows from the Review table.
package srs

import (
	"time"
)

const day time.Duration = 86400000000000 // In nanoseconds

type Review struct {
	Due      time.Time
	Interval time.Duration
	Reviewed time.Time
	Correct  bool
}

// Computes next review schedule.
// If review is nil, creates Review with default values for initial review.
func nextReview(review *Review, correct bool) Review {
	var interval time.Duration = 0
	if correct {
		if review != nil {
			interval = 2 * review.Interval
		} else {
			interval = day
		}
	}

	now := time.Now()
	return Review{
		Reviewed: now,
		Interval: interval,
		Due:      now.Add(interval),
		Correct:  correct,
	}
}
