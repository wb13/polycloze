package srs

import (
	"database/sql"
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

// Schedule due words, no more than count.
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

// Get most recent review of word.
// Result is nil whenever something goes wrong.
func mostRecentReview(tx *sql.Tx, word string) *Review {
	query := `
SELECT due, interval, reviewed, correct FROM MostRecentReview
WHERE word = ?
`
	var review Review
	row := tx.QueryRow(query, word)
	err := row.Scan(
		&review.Due,
		&review.Interval,
		&review.Reviewed,
		&review.Correct,
	)
	if err != nil {
		return nil
	}
	return &review
}

// Update review status of word.
func (ws *WordScheduler) Update(word string, correct bool) error {
	tx, err := ws.db.Begin()
	if err != nil {
		return err
	}

	review := mostRecentReview(tx, word)

	query := `
INSERT INTO Review (word, interval, due, correct)
VALUES (?, ?, ?, ?)
`
	var next Review
	if review == nil {
		next = defaultReview(correct)
	} else {
		next = nextReview(review, correct)
	}

	_, err = tx.Exec(query, word, next.Interval, next.Due, correct)
	if err != nil {
		return err
	}
	return tx.Commit()
}
