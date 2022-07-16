// Combines word_queue and review_scheduler to schedule words.
package word_scheduler

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lggruspe/polycloze/database"
	rs "github.com/lggruspe/polycloze/review_scheduler"
	"github.com/lggruspe/polycloze/text"
	"github.com/lggruspe/polycloze/word_queue"
)

// Gets preferred difficulty/frequency_class.
func preferredDifficulty(s *database.Session) int {
	query := `select frequency_class from student`

	var difficulty int
	row := s.QueryRow(query)
	row.Scan(&difficulty)
	return difficulty
}

// NOTE Only returns new words if words for review < n.
func GetWords(s *database.Session, n int) ([]string, error) {
	reviews, err := rs.ScheduleReviewNow(s, n)
	if err != nil {
		return nil, err
	}
	words, err := word_queue.GetNewWords(s, n-len(reviews), preferredDifficulty(s))
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
	words, err := word_queue.GetNewWords(s, n-len(reviews), preferredDifficulty(s))
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
	words, err := word_queue.GetNewWordsWith(s, n-len(reviews), preferredDifficulty(s), pred)
	if err != nil {
		return nil, err
	}
	return append(reviews, words[:]...), nil
}

func frequencyClass(s *database.Session, word string) int {
	query := `select frequency_class from l2.word where word = ?`
	row := s.QueryRow(query, text.Casefold(word))

	var frequency_class int
	row.Scan(&frequency_class)
	return frequency_class
}

func isNewWord(s *database.Session, word string) bool {
	query := `select rowid from review where item = ?`
	row := s.QueryRow(query, text.Casefold(word))

	var rowid int
	err := row.Scan(&rowid)
	return err != nil && errors.Is(err, sql.ErrNoRows)
}

// This should only be called when an item is seen for the first time.
func updateStudentStats(s *database.Session, correct bool) error {
	query := `update student set correct = correct + 1`
	if !correct {
		query = `update student set incorrect = incorrect + 1`
	}
	_, err := s.Exec(query)
	return err
}

func UpdateWord(s *database.Session, word string, correct bool) error {
	if frequencyClass(s, word) >= preferredDifficulty(s) && isNewWord(s, word) {
		if err := updateStudentStats(s, correct); err != nil {
			return err
		}
	}
	return rs.UpdateReview(s, text.Casefold(word), correct)
}

// See UpdateReviewAt.
func UpdateWordAt(s *database.Session, word string, correct bool, at time.Time) error {
	if frequencyClass(s, word) >= preferredDifficulty(s) && isNewWord(s, word) {
		if err := updateStudentStats(s, correct); err != nil {
			return err
		}
	}
	return rs.UpdateReviewAt(s, text.Casefold(word), correct, at)
}
