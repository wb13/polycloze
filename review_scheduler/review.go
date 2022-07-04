// Defines the Review struct for storing rows from the Review table.
package review_scheduler

import (
	"database/sql"
	"time"
)

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
	now := time.Now().UTC()
	reviewed := now
	if review != nil {
		if now.Before(review.Due) {
			// Don't increase interval if the user crammed
			return review.Interval, nil
		}
		reviewed = review.Reviewed
	}

	interval := now.Sub(reviewed) // this is greater than review.Interval
	return nextInterval(tx, interval)
}

func updateIntervalStats(tx *sql.Tx, review *Review, correct bool) error {
	var interval time.Duration = 0
	if review != nil {
		interval = review.Interval
	}

	query := `update interval set correct = correct + 1 where interval = ?`
	if !correct {
		query = `update interval set incorrect = incorrect + 1 where interval = ?`
	}
	_, err := tx.Exec(query, interval)
	return err
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
