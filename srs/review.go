// Defines the Review struct for storing rows from the Review table.
package srs

import (
	"math"
	"time"
)

const day time.Duration = 86400000000000 // In nanoseconds

type Review struct {
	Due      time.Time     // Due date of next review
	Interval time.Duration // Interval between now and due date
	Reviewed time.Time
	Correct  bool
}

func (r Review) Level() int {
	return int(math.Floor(math.Log2(2*r.Interval.Hours()/24 + 1)))
}

// Returns computed level of review item.
func getLevel(r *Review) int {
	if r == nil {
		return 0
	}
	return r.Level()
}

// Computes next review schedule.
// If review is nil, creates Review with default values for initial review.
func nextReview(review *Review, correct bool, coefficient float64) Review {
	var interval time.Duration = 0
	if correct {
		if review != nil {
			interval = time.Duration(coefficient * float64(review.Interval.Nanoseconds()))
		} else {
			interval = day
		}
	}

	now := time.Now().UTC()
	return Review{
		Reviewed: now,
		Interval: interval,
		Due:      now.Add(interval),
		Correct:  correct,
	}
}
