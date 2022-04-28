package srs

import (
	"database/sql"
	"embed"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var fs embed.FS

var day time.Duration

func init() {
	day, _ = time.ParseDuration("24h")
}

// Migrate to latest version of the database.
func migrateUp(db *sql.DB) error {
	dbDriver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return err
	}

	srcDriver, err := iofs.New(fs, "migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance(
		"iofs",
		srcDriver,
		"sqlite3",
		dbDriver,
	)
	if err != nil {
		return err
	}

	return m.Up()
}

// type SpacingAlgorithm func(*WordScheduler, Review, bool) (Review, error) // TODO

type WordScheduler struct {
	db *sql.DB
	// algorithm SpacingAlgorithm	// TODO
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

type Review struct {
	Word     string
	Due      time.Time
	Interval time.Duration
	Reviewed time.Time
	Correct  bool
}

// Value of first-time word review
func defaultReview(word string, correct bool) Review {
	interval := day
	if !correct {
		interval = 0
	}

	now := time.Now()
	return Review{
		Word:     word,
		Reviewed: now,
		Interval: interval,
		Due:      now.Add(interval),
		Correct:  correct,
	}
}

// Get most recent review of word.
// Result is nil whenever something goes wrong.
func mostRecentReview(tx *sql.Tx, word string) *Review {
	query := `
SELECT word, due, interval, reviewed, correct FROM MostRecentReview
WHERE word = ?
`
	var review Review
	row := tx.QueryRow(query, word)
	err := row.Scan(
		&review.Word,
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

// Computes next review schedule.
func nextReview(review *Review, correct bool) Review {
	due := review.Due
	interval := review.Interval

	if correct {
		interval *= 2
	} else {
		interval = 0
	}
	due = due.Add(interval)

	return Review{
		Word:     review.Word,
		Due:      due,
		Interval: interval,
		Reviewed: time.Now(),
		Correct:  correct,
	}
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
		next = defaultReview(word, correct)
	} else {
		next = nextReview(review, correct)
	}

	_, err = tx.Exec(query, word, next.Interval, next.Due, correct)
	if err != nil {
		return err
	}
	return tx.Commit()
}
