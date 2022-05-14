// Combines word_queue and review_scheduler to schedule words.
package word_scheduler

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"github.com/lggruspe/polycloze/database"
	rs "github.com/lggruspe/polycloze/review_scheduler"
	"github.com/lggruspe/polycloze/word_queue"
)

// Returns in-memory sqlite DB, and attaches specified databases.
func New(reviewDB, l2db string) (*sql.DB, error) {
	if err := database.UpgradeFile(reviewDB); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	if err := database.Attach(db, "review", reviewDB); err != nil {
		return nil, err
	}
	if err := database.Attach(db, "l2", l2db); err != nil {
		return nil, err
	}
	return db, nil
}

// NOTE Only returns new words if words for review < n.
func GetWords(s *database.Session, n int) ([]string, error) {
	reviews, err := rs.ScheduleReviewNow(s, n)
	if err != nil {
		return nil, err
	}
	words, err := word_queue.GetNewWords(s, n-len(reviews))
	if err != nil {
		return nil, err
	}
	return append(reviews, words[:]...), nil
}
