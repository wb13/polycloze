// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package word_scheduler

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lggruspe/polycloze/database"
	rs "github.com/lggruspe/polycloze/review_scheduler"
	"github.com/lggruspe/polycloze/text"
)

// Gets preferred difficulty/frequency_class.
func PreferredDifficulty(s *database.Session) int {
	query := `select frequency_class from student`

	var difficulty int
	row := s.QueryRow(query)
	_ = row.Scan(&difficulty)
	return difficulty
}

// Same as GetWords, but takes an additional time.Time argument.
func GetWordsAt(s *database.Session, n int, due time.Time) ([]string, error) {
	reviews, err := rs.ScheduleReview(s, due, n)
	if err != nil {
		return nil, err
	}
	words, err := GetNewWordsWith(s, n-len(reviews), PreferredDifficulty(s), func(_ string) bool {
		return true
	})
	if err != nil {
		return nil, err
	}
	return append(reviews, words...), nil
}

// Returns up to words to make flashcards for.
// Only includes words that satisfy the predicate.
func GetWordsWith(s *database.Session, n int, pred func(word string) bool) ([]string, error) {
	reviews, err := rs.ScheduleReviewNowWith(s, n, pred)
	if err != nil {
		return nil, err
	}
	words, err := GetNewWordsWith(s, n-len(reviews), PreferredDifficulty(s), pred)
	if err != nil {
		return nil, err
	}
	return append(reviews, words...), nil
}

func frequencyClass(s *database.Session, word string) int {
	query := `select frequency_class from word where word = ?`
	row := s.QueryRow(query, text.Casefold(word))

	var result int
	_ = row.Scan(&result)
	return result
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
	if frequencyClass(s, word) >= PreferredDifficulty(s) && isNewWord(s, word) {
		if err := updateStudentStats(s, correct); err != nil {
			return err
		}
	}
	if err := rs.UpdateReview(s, text.Casefold(word), correct); err != nil {
		return err
	}
	return postTune(s)
}

// See UpdateReviewAt.
func UpdateWordAt(s *database.Session, word string, correct bool, at time.Time) error {
	if frequencyClass(s, word) >= PreferredDifficulty(s) && isNewWord(s, word) {
		if err := updateStudentStats(s, correct); err != nil {
			return err
		}
	}
	if err := rs.UpdateReviewAt(s, text.Casefold(word), correct, at); err != nil {
		return err
	}
	return postTune(s)
}
