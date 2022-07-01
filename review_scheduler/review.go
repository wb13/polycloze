// Defines the Review struct for storing rows from the Review table.
package review_scheduler

import (
	"database/sql"
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

// Calculates interval for next review.
func calculateInterval(tx *sql.Tx, review *Review, correct bool) (time.Duration, error) {
	if !correct {
		return 0, nil
	}
	if review == nil || review.Interval == 0 {
		return day, nil
	}

	now := time.Now().UTC()
	if now.Before(review.Due) {
		// Don't increase interval if the user crammed
		return review.Interval, nil
	}

	interval := now.Sub(review.Reviewed)
	return nextInterval(tx, interval)
}

// Computes next review schedule.
// If review is nil, creates Review with default values for initial review.
func nextReview(tx *sql.Tx, review *Review, correct bool) (Review, error) {
	var r Review
	interval, err := calculateInterval(tx, review, correct)
	if err != nil {
		return r, err
	}
	now := time.Now().UTC()
	r.Reviewed = now
	r.Interval = interval
	r.Due = now.Add(interval)
	return r, nil
}
