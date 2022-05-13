// Combines word_queue and review_scheduler to schedule words.
package word_scheduler

import (
	"database/sql"
	"path"

	_ "github.com/mattn/go-sqlite3"

	"github.com/lggruspe/polycloze-srs/database"
	rs "github.com/lggruspe/polycloze-srs/review_scheduler"
	"github.com/lggruspe/polycloze-srs/word_queue"
)

// Returns in-memory sqlite DB, and attaches specified databases.
func New(reviewDB, l2db string) (*sql.DB, error) {
	migrations := path.Join("migrations", "review_scheduler")
	if err := database.UpgradeFile(reviewDB, migrations); err != nil {
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
// Expects language DB and review DB to already be attached.
func GetWords(db *sql.DB, n int) ([]string, error) {
	reviews, err := rs.ScheduleReviewNow(db, n)
	if err != nil {
		return nil, err
	}
	words, err := word_queue.GetNewWords(db, n-len(reviews))
	if err != nil {
		return nil, err
	}
	return append(reviews, words[:]...), nil
}
