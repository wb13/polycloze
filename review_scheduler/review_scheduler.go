// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package review_scheduler

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lggruspe/polycloze/database"
	_ "github.com/mattn/go-sqlite3"
)

// Returns items due for review, no more than count.
// Pass a negative count if you want to get all due items.
func ScheduleReview[T database.Querier](q T, due time.Time, count int) ([]string, error) {
	query := `SELECT item FROM review WHERE due <= ? ORDER BY due LIMIT ?`
	rows, err := q.Query(query, due.Unix(), count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []string
	for rows.Next() {
		var item string
		if err := rows.Scan(&item); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// Same as Schedule, but with some default args.
func ScheduleReviewNow[T database.Querier](q T, count int) ([]string, error) {
	return ScheduleReview(q, time.Now(), count)
}

// Same as ScheduleReviewNowWith, but takes a predicate argument.
// Only items that satisfy the predicate are included in the result.
func ScheduleReviewNowWith[T database.Querier](q T, count int, pred func(item string) bool) ([]string, error) {
	query := `select item from review where due < ? order by due`
	rows, err := q.Query(query, time.Now().UTC())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []string
	for rows.Next() && len(items) < count {
		var item string
		if err := rows.Scan(&item); err != nil {
			return nil, err
		}
		if pred(item) {
			items = append(items, item)
		}
	}
	return items, nil
}

// Gets most recent review of item.
func mostRecentReview(tx *sql.Tx, item string) (*Review, error) {
	query := `select due, interval, reviewed from review where item = ?`
	row := tx.QueryRow(query, item)
	var review Review

	var interval time.Duration
	var due, reviewed int64
	err := row.Scan(
		&due,
		&interval,
		&reviewed,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	review.Due = time.Unix(due, 0)
	review.Reviewed = time.Unix(reviewed, 0)
	review.Interval = interval * time.Second
	return &review, nil
}

// Updates review status of item.
func UpdateReviewAt[T database.Querier](q T, item string, correct bool, now time.Time) error {
	tx, err := q.Begin()
	if err != nil {
		return err
	}

	review, err := mostRecentReview(tx, item)
	if err != nil {
		return err
	}

	if review == nil || !now.Before(review.Due) {
		// Only update interval stats if the student didn't cram
		if err := updateIntervalStats(tx, review, correct); err != nil {
			return err
		}
	}

	next, err := nextReview(tx, review, correct, now)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO review (item, interval, due, learned, reviewed)
		VALUES (?, ?, ?, unixepoch('now'), unixepoch('now'))
		ON CONFLICT (item) DO UPDATE SET
			interval = excluded.interval,
			due = excluded.due,
			reviewed = excluded.reviewed
	`
	_, err = tx.Exec(
		query,
		item,
		seconds(next.Interval),
		next.Due.Unix(),
	)
	if err != nil {
		return err
	}
	if err := autoTune(tx); err != nil {
		return err
	}
	return tx.Commit()
}

func UpdateReview[T database.Querier](q T, item string, correct bool) error {
	return UpdateReviewAt(q, item, correct, time.Now().UTC())
}
