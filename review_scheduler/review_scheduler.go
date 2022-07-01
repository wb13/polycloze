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

// Returns items due for review, no more than count.
// Pass a negative count if you want to get all due items.
func ScheduleReview(s *database.Session, due time.Time, count int) ([]string, error) {
	query := `select item from review where due < ? order by due limit ?`
	rows, err := s.Query(query, due.UTC(), count)
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
func ScheduleReviewNow(s *database.Session, count int) ([]string, error) {
	return ScheduleReview(s, time.Now().UTC(), count)
}

// Same as ScheduleReviewNowWith, but takes a predicate argument.
// Only items that satisfy the predicate are included in the result.
func ScheduleReviewNowWith(s *database.Session, count int, pred func(item string) bool) ([]string, error) {
	query := `select item from review where due < ? order by due`
	rows, err := s.Query(query, time.Now().UTC())
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
func UpdateReview(s *database.Session, item string, correct bool) error {
	tx, err := s.Begin()
	if err != nil {
		return err
	}

	if err := initTuner(tx); err != nil {
		return err
	}

	review, err := mostRecentReview(tx, item)
	if err != nil {
		return err
	}

	query := `
insert into review (item, interval, due) values (?, ?, ?)
	on conflict (item) do update set interval=excluded.interval, due=excluded.due
`

	next, err := nextReview(tx, review, correct)
	if err != nil {
		return err
	}

	_, err = tx.Exec(query, item, next.Interval, next.Due)
	if err != nil {
		return err
	}
	if err := autoTune(tx); err != nil {
		return err
	}
	return tx.Commit()
}
