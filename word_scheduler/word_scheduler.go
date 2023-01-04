// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package word_scheduler

import (
	"time"

	"github.com/polycloze/polycloze/database"
	"github.com/polycloze/polycloze/difficulty"
	rs "github.com/polycloze/polycloze/review_scheduler"
	"github.com/polycloze/polycloze/text"
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

	level := difficulty.GetLatest(q).Level
	words, err := GetNewWordsWith(q, n-len(reviews), level, func(_ string) bool {
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

	level := difficulty.GetLatest(q).Level
	words, err := GetNewWordsWith(q, n-len(reviews), level, pred)
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

func UpdateWord[T database.Querier](q T, word string, correct bool) error {
	return rs.UpdateReview(q, text.Casefold(word), correct)
}

// See UpdateReviewAt.
func UpdateWordAt[T database.Querier](q T, word string, correct bool, at time.Time) error {
	return rs.UpdateReviewAt(q, text.Casefold(word), correct, at)
}

type ReviewResult = rs.Result

// Saves word review results in bulk.
func BulkSaveWords[T database.Querier](q T, reviews []ReviewResult, at time.Time) error {
	// Client already casefolds words, but let's casefold again to be sure.
	for i, review := range reviews {
		reviews[i].Word = text.Casefold(review.Word)
	}
	return rs.BulkSaveReviews(q, reviews, at)
}
