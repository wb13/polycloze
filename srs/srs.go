package srs

import (
	"database/sql"
	"errors"
	"time"
)

type WordScheduler struct {
	db *sql.DB
}

func InitWordScheduler(db *sql.DB) (WordScheduler, error) {
	if err := migrateUp(db); err != nil {
		return WordScheduler{nil}, err
	}
	return WordScheduler{db}, nil
}

// Returns words that are due for review, no more than count.
// Pass a negative count if you want to get all due words.
func (ws *WordScheduler) Schedule(due time.Time, count int) ([]string, error) {
	query := `
SELECT word FROM MostRecentReview WHERE due < ? LIMIT ?
`
	rows, err := ws.db.Query(query, due, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []string
	for rows.Next() {
		var word string
		if err := rows.Scan(&word); err != nil {
			return nil, err
		}
		words = append(words, word)
	}
	return words, nil
}

// Same as Schedule, but with some default args.
func (ws *WordScheduler) ScheduleNow(count int) ([]string, error) {
	return ws.Schedule(time.Now(), count)
}

// Gets most recent review of word.
func mostRecentReview(tx *sql.Tx, word string) (*Review, error) {
	query := `
SELECT due, interval, reviewed, correct, streak FROM MostRecentReview
WHERE word = ?
`
	row := tx.QueryRow(query, word)
	var review Review

	var due string
	var reviewed string
	err := row.Scan(
		&due,
		&review.Interval,
		&reviewed,
		&review.Correct,
		&review.Streak,
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

// Updates review status of word.
func (ws *WordScheduler) Update(word string, correct bool) error {
	tx, err := ws.db.Begin()
	if err != nil {
		return err
	}

	review, err := mostRecentReview(tx, word)
	if err != nil {
		return err
	}

	query := `
INSERT INTO Review (word, interval, due, correct, streak)
VALUES (?, ?, ?, ?, ?)
`

	coefficient := getCoefficient(tx, getStreak(review))
	next := nextReview(review, correct, coefficient)
	_, err = tx.Exec(query, word, next.Interval, next.Due, correct, next.Streak)
	if err != nil {
		return err
	}
	if err := autoTune(tx); err != nil {
		return err
	}
	return tx.Commit()
}
