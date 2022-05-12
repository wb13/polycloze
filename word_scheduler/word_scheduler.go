// Combines word_queue and review_scheduler to schedule words.
package word_scheduler

import (
	"database/sql"

	rs "github.com/lggruspe/polycloze-srs/review_scheduler"
	"github.com/lggruspe/polycloze-srs/word_queue"
)

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
