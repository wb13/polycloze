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

// Same as GetWords, but takes an additional time.Time argument.
func GetWordsAt[T database.Querier](q T, n int, due time.Time) ([]Word, error) {
	var result []Word
	reviews, err := rs.ScheduleReview(q, due, n)
	if err != nil {
		return nil, err
	}
	for _, word := range reviews {
		result = append(result, Word{
			Word: word,
			New:  false,
		})
	}

	words, err := GetNewWordsWith(q, n-len(reviews), Placement(q), func(_ string) bool {
		return true
	})
	if err != nil {
		return nil, err
	}
	return append(result, words...), nil
}

// Returns up to words to make flashcards for.
// Only includes words that satisfy the predicate.
func GetWordsWith[T database.Querier](q T, n int, pred func(word string) bool) ([]Word, error) {
	var result []Word

	reviews, err := rs.ScheduleReviewNowWith(q, n, pred)
	if err != nil {
		return nil, err
	}
	for _, word := range reviews {
		result = append(result, Word{
			Word: word,
			New:  false,
		})
	}

	words, err := GetNewWordsWith(q, n-len(reviews), Placement(q), pred)
	if err != nil {
		return nil, err
	}
	return append(result, words...), nil
}

func frequencyClass[T database.Querier](q T, word string) int {
	query := `select frequency_class from word where word = ?`
	row := q.QueryRow(query, text.Casefold(word))

	var result int
	_ = row.Scan(&result)
	return result
}

func isNewWord[T database.Querier](q T, word string) bool {
	query := `select rowid from review where item = ?`
	row := q.QueryRow(query, text.Casefold(word))

	var rowid int
	err := row.Scan(&rowid)
	return err != nil && errors.Is(err, sql.ErrNoRows)
}

func UpdateWord[T database.Querier](q T, word string, correct bool) error {
	if isNewWord(q, word) {
		class := frequencyClass(q, word)
		if err := updateNewWordStat(q, class, correct); err != nil {
			return err
		}
	}
	return rs.UpdateReview(q, text.Casefold(word), correct)
}

// See UpdateReviewAt.
func UpdateWordAt[T database.Querier](q T, word string, correct bool, at time.Time) error {
	if isNewWord(q, word) {
		class := frequencyClass(q, word)
		if err := updateNewWordStat(q, class, correct); err != nil {
			return err
		}
	}
	return rs.UpdateReviewAt(q, text.Casefold(word), correct, at)
}
