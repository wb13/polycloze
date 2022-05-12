// Combines word_queue and review_scheduler to schedule words.
package word_scheduler

import (
	"database/sql"

	_ "github.com/lggruspe/polycloze-srs/review_scheduler"
	"github.com/lggruspe/polycloze-srs/word_queue"
)

// NOTE Only returns new words if words for review < n.
// Expects language DB and review DB to already be attached.
func GetWords(db *sql.DB, n int) ([]string, error) {
	return word_queue.GetNewWords(db, n)
}
