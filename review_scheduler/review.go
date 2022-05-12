// Defines the Review struct for storing rows from the Review table.
package review_scheduler

import (
	"math"
	"time"
)

const day time.Duration = 86400000000000 // In nanoseconds

type Review struct {
	Due      time.Time     // Due date of next review
	Interval time.Duration // Interval between now and due date
	Reviewed time.Time
}

func (r Review) Correct() bool {
	return r.Interval > 0
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

// Calculates interval for next review.
func calculateInterval(review *Review, correct bool, coefficient float64) time.Duration {
	if !correct {
		return 0
	}
	if review == nil {
		return day
	}
	if review.Interval == 0 {
		return day
	}

	now := time.Now().UTC()
	if now.Before(review.Due) {
		return review.Interval
	}

	interval := now.Sub(review.Reviewed)
	return time.Duration(coefficient * float64(interval.Nanoseconds()))
}

// Computes next review schedule.
// If review is nil, creates Review with default values for initial review.
func nextReview(review *Review, correct bool, coefficient float64) Review {
	interval := calculateInterval(review, correct, coefficient)
	now := time.Now().UTC()
	return Review{
		Reviewed: now,
		Interval: interval,
		Due:      now.Add(interval),
	}
}
