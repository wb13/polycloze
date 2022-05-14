package review_scheduler

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/lggruspe/polycloze/database"
)

// Returns sql.DB with review_scheduler schema.
func New(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	if err := database.Upgrade(db); err != nil {
		return nil, err
	}
	return db, nil
}

func InitAttached(db *sql.DB, dbPath string) error {
	if err := database.UpgradeFile(dbPath); err != nil {
		return err
	}
	return database.Attach(db, "review", dbPath)
}

// Returns items due for review, no more than count.
// Pass a negative count if you want to get all due items.
//
// Expects db to contain updated review_scheduler schema.
func ScheduleReview(db *sql.DB, due time.Time, count int) ([]string, error) {
	query := `SELECT item FROM most_recent_review WHERE due < ? LIMIT ?`
	rows, err := db.Query(query, due.UTC(), count)
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
func ScheduleReviewNow(db *sql.DB, count int) ([]string, error) {
	return ScheduleReview(db, time.Now().UTC(), count)
}

// Gets most recent review of item.
func mostRecentReview(tx *sql.Tx, item string) (*Review, error) {
	query := `
SELECT due, interval, reviewed FROM most_recent_review
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
func UpdateReview(db *sql.DB, item string, correct bool) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	review, err := mostRecentReview(tx, item)
	if err != nil {
		return err
	}

	query := `
INSERT INTO review (item, interval, due)
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
