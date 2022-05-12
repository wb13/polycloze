// Returns words that the user hasn't learned.
package word_queue

import (
	"database/sql"
)

// Gets up to n new words from db.
// Pass a negative n if you don't want a word limit.
// Expects language DB and review DB to be attached.
func GetNewWord(db *sql.DB, n int) ([]string, error) {
	query := `
select word from word where word not in
(select item from review)
order by frequency desc
limit ?
`
	rows, err := db.Query(query, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []string
	for rows.Next() {
		var word string
		if err := rows.Scan(&word); err != nil {
			return nil, err
		}
		words = append(words, word)
	}
	return words, nil
}
