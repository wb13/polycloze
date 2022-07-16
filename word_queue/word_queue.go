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

// Gets up to n new words from db.
// Pass a negative n if you don't want a word limit.
func GetNewWords(s *database.Session, n int) ([]string, error) {
	query := `
select word from l2.word where word not in
(select item from review)
order by frequency desc
limit ?
`
	rows, err := s.Query(query, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return getNRows(rows, n, func(_ string) bool {
		return true
	})
}

// Same as GetNewWords, but takes a predicate argument.
// Only words that satisfy the predicate are included in the result.
func GetNewWordsWith(s *database.Session, n int, pred func(word string) bool) ([]string, error) {
	query := `
select word from l2.word where word not in
(select item from review)
order by frequency desc
`
	rows, err := s.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return getNRows(rows, n, pred)
}
