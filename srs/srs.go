package srs

import (
	"database/sql"
	"errors"
	"time"
)

type ReviewScheduler struct {
	db *sql.DB
}

func InitReviewScheduler(db *sql.DB) (ReviewScheduler, error) {
	if err := migrateUp(db); err != nil {
		return ReviewScheduler{nil}, err
	}
	return ReviewScheduler{db}, nil
}

// Returns items due for review, no more than count.
// Pass a negative count if you want to get all due items.
func (ws *ReviewScheduler) Schedule(due time.Time, count int) ([]string, error) {
	query := `
SELECT item FROM MostRecentReview WHERE due < ? LIMIT ?
`
	rows, err := ws.db.Query(query, due.UTC(), count)
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
func (ws *ReviewScheduler) ScheduleNow(count int) ([]string, error) {
	return ws.Schedule(time.Now().UTC(), count)
}

// Gets most recent review of item.
func mostRecentReview(tx *sql.Tx, item string) (*Review, error) {
	query := `
SELECT due, interval, reviewed FROM MostRecentReview
WHERE item = ?
`
	row := tx.QueryRow(query, item)
	var review Review

	var due string
	var reviewed string
	err := row.Scan(
		&due,
		&review.Interval,
		&reviewed,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	parsedDue, err := parseTimestamp(due)
	if err != nil {
		return nil, err
	}
	parsedReviewed, err := parseTimestamp(reviewed)
	if err != nil {
		return nil, err
	}
	review.Due = parsedDue
	review.Reviewed = parsedReviewed
	return &review, nil
}

// Updates review status of item.
func (ws *ReviewScheduler) Update(item string, correct bool) error {
	tx, err := ws.db.Begin()
	if err != nil {
		return err
	}

	review, err := mostRecentReview(tx, item)
	if err != nil {
		return err
	}

	query := `
INSERT INTO Review (item, interval, due)
VALUES (?, ?, ?)
`

	coefficient := getCoefficient(tx, getLevel(review))
	next := nextReview(review, correct, coefficient)
	_, err = tx.Exec(query, item, next.Interval, next.Due)
	if err != nil {
		return err
	}
	if err := autoTune(tx); err != nil {
		return err
	}
	return tx.Commit()
}
