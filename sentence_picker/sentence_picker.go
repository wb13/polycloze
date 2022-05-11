package sentence_picker

import (
	"database/sql"
	"encoding/json"
	"errors"
)

func InitSentencePicker(db *sql.DB, langDB, reviewDB string) error {
	if err := migrateUp(db); err != nil {
		return err
	}
	if err := attachDatabase(db, "language_schema", langDB); err != nil {
		return err
	}
	if err := attachDatabase(db, "review_schema", reviewDB); err != nil {
		return err
	}
	return nil
}

// Returns "easiest" sentence that contains word.
func PickSentence(db *sql.DB, word string) ([]string, error) {
	query := `
select tokens, min(difficulty)
from sentence
join contains on (sentence.id = contains.sentence)
join word on (word.id = contains.word)
join sentence_difficulty using (sentence)
where word.word = ?
`
	row := db.QueryRow(query, word)

	var jsonStr string
	var difficulty float64
	if err := row.Scan(&jsonStr, &difficulty); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		// doesn't panic, because words in db are guaranteed to be in a sentence
		panic("no sentences found")
	}

	var tokens []string
	if err := json.Unmarshal([]byte(jsonStr), &tokens); err != nil {
		return nil, err
	}
	return tokens, nil
}
