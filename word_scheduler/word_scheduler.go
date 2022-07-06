// Combines word_queue and review_scheduler to schedule words.
package word_scheduler

import (
	"time"

	"github.com/lggruspe/polycloze/database"
	rs "github.com/lggruspe/polycloze/review_scheduler"
	"github.com/lggruspe/polycloze/text"
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

// Same as GetWords, but takes an additional time.Time argument.
func GetWordsAt(s *database.Session, n int, due time.Time) ([]string, error) {
	reviews, err := rs.ScheduleReview(s, due, n)
	if err != nil {
		return nil, err
	}
	words, err := word_queue.GetNewWords(s, n-len(reviews))
	if err != nil {
		return nil, err
	}
	return append(reviews, words[:]...), nil
}

// Same as GetWords, but takes an additional predicate argument.
// Only includes words that satisfy the predicate.
func GetWordsWith(s *database.Session, n int, pred func(word string) bool) ([]string, error) {
	reviews, err := rs.ScheduleReviewNowWith(s, n, pred)
	if err != nil {
		return nil, err
	}
	words, err := word_queue.GetNewWordsWith(s, n-len(reviews), pred)
	if err != nil {
		return nil, err
	}
	return append(reviews, words[:]...), nil
}

func UpdateWord(s *database.Session, word string, correct bool) error {
	return rs.UpdateReview(s, text.Casefold(word), correct)
}

// See UpdateReviewAt.
func UpdateWordAt(s *database.Session, word string, correct bool, at time.Time) error {
	return rs.UpdateReviewAt(s, text.Casefold(word), correct, at)
}
