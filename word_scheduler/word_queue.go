// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

// Returns words that the user hasn't learned.
package word_scheduler

import (
	"database/sql"

	"github.com/lggruspe/polycloze/database"
)

// Scans rows in query (each row has `word` and `frequency_class`).
// Does not close Rows.
func getNRows(rows *sql.Rows, n int, pred func(word string) bool) ([]Word, error) {
	var words []Word
	for rows.Next() && len(words) < n {
		var word string
		var frequencyClass int
		if err := rows.Scan(&word, &frequencyClass); err != nil {
			return nil, err
		}
		if pred(word) {
			words = append(words, Word{
				Word:       word,
				New:        true,
				Difficulty: frequencyClass,
			})
		}
	}
	return words, nil
}

func getWordsAboveDifficultyWith[T database.Querier](q T, n, preferredDifficulty int, pred func(word string) bool) ([]Word, error) {
	query := `
		SELECT word, frequency_class
		FROM word
		WHERE frequency_class >= ? AND word NOT IN (
			SELECT item FROM review
		)
		ORDER BY id ASC
`
	rows, err := q.Query(query, preferredDifficulty)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return getNRows(rows, n, pred)
}

func getWordsBelowDifficultyWith[T database.Querier](q T, n, preferredDifficulty int, pred func(word string) bool) ([]Word, error) {
	query := `
		SELECT word, frequency_class
		FROM word
		WHERE frequency_class < ? AND word NOT IN (
			SELECT item FROM review
		)
		ORDER BY id DESC
`
	rows, err := q.Query(query, preferredDifficulty)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return getNRows(rows, n, pred)
}

// Gets up to n new words from db.
// Pass a negative n if you don't want a word limit.
// Uses preferredDifficulty as minimum word frequency class.
// If there are not enough words in query result, will also include words below
// the preferredDifficulty.
// Only words that satisfy the predicate are included in the result.
func GetNewWordsWith[T database.Querier](q T, n, preferredDifficulty int, pred func(word string) bool) ([]Word, error) {
	words, err := getWordsAboveDifficultyWith(q, n, preferredDifficulty, pred)
	if err != nil {
		return nil, err
	}
	if preferredDifficulty <= 0 || len(words) >= n {
		return words, nil
	}

	more, err := getWordsBelowDifficultyWith(q, n-len(words), preferredDifficulty, pred)
	if err != nil {
		return nil, err
	}
	return append(words, more...), nil
}
