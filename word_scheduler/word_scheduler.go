// Combines word_queue and review_scheduler to schedule words.
package word_scheduler

import (
	"github.com/lggruspe/polycloze/database"
	rs "github.com/lggruspe/polycloze/review_scheduler"
	"github.com/lggruspe/polycloze/word_queue"
)

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
