// Returns words that the user hasn't learned.
package word_queue

import (
	"database/sql"

	"github.com/lggruspe/polycloze/database"
)

// NOTE Does not close Rows.
func getNRows(rows *sql.Rows, n int, pred func(word string) bool) ([]string, error) {
	var words []string
	for rows.Next() && len(words) < n {
		var word string
		if err := rows.Scan(&word); err != nil {
			return nil, err
		}
		if pred(word) {
			words = append(words, word)
		}
	}
	return words, nil
}

func getWordsAboveDifficulty(s *database.Session, n, preferredDifficulty int) ([]string, error) {
	query := `
select word from word where frequency_class >= ? and word not in
(select item from review)
order by id asc limit ?
`
	rows, err := s.Query(query, preferredDifficulty, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return getNRows(rows, n, func(_ string) bool {
		return true
	})
}

func getWordsBelowDifficulty(s *database.Session, n, preferredDifficulty int) ([]string, error) {
	query := `
select word from word where frequency_class < ? and word not in
(select item from review)
order by id desc limit ?
`
	rows, err := s.Query(query, preferredDifficulty, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return getNRows(rows, n, func(_ string) bool {
		return true
	})
}

func getWordsAboveDifficultyWith(s *database.Session, n, preferredDifficulty int, pred func(word string) bool) ([]string, error) {
	query := `
select word from word where frequency_class >= ? and word not in
(select item from review)
order by id asc
`
	rows, err := s.Query(query, preferredDifficulty, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return getNRows(rows, n, pred)
}

func getWordsBelowDifficultyWith(s *database.Session, n, preferredDifficulty int, pred func(word string) bool) ([]string, error) {
	query := `
select word from word where frequency_class < ? and word not in
(select item from review)
order by id desc
`
	rows, err := s.Query(query, preferredDifficulty, n)
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
func GetNewWords(s *database.Session, n, preferredDifficulty int) ([]string, error) {
	words, err := getWordsAboveDifficulty(s, n, preferredDifficulty)
	if err != nil {
		return nil, err
	}
	if preferredDifficulty <= 0 || len(words) >= n {
		return words, nil
	}

	more, err := getWordsBelowDifficulty(s, n-len(words), preferredDifficulty)
	if err != nil {
		return nil, err
	}
	words = append(words, more...)
	return words, nil
}

// Same as GetNewWords, but takes a predicate argument.
// Only words that satisfy the predicate are included in the result.
func GetNewWordsWith(s *database.Session, n, preferredDifficulty int, pred func(word string) bool) ([]string, error) {
	words, err := getWordsAboveDifficultyWith(s, n, preferredDifficulty, pred)
	if err != nil {
		return nil, err
	}
	if preferredDifficulty <= 0 || len(words) >= n {
		return words, nil
	}

	more, err := getWordsBelowDifficultyWith(s, n-len(words), preferredDifficulty, pred)
	if err != nil {
		return nil, err
	}
	words = append(words, more...)
	return words, nil
}
