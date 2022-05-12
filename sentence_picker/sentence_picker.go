package sentence_picker

import (
	"database/sql"
	"encoding/json"
	"math"
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

// Returns map of sentence ID -> sentence difficulty of sentences containing word.
func sentencesThatContain(db *sql.DB, word string) (map[int]float64, error) {
	// NOTE seen counter is only used to add variation in the returned sentences
	query := `
select contains.sentence as sentence,
			count(word) + sum(difficulty) + coalesce(counter, 0.0) as difficulty
from contains
join word_difficulty using(word)
left join seen on (contains.sentence = seen.sentence)
where contains.sentence in
	(select sentence from contains join word on (contains.word = word.id) where word.word = ?)
group by contains.sentence;
`
	rows, err := db.Query(query, word)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sentences := make(map[int]float64)
	for rows.Next() {
		var id int
		var difficulty float64
		if err := rows.Scan(&id, &difficulty); err != nil {
			return nil, err
		}
		sentences[id] = difficulty
	}
	return sentences, nil
}

func sentenceTokens(db *sql.DB, sentence int) ([]string, error) {
	query := `select tokens from sentence where id = ?`
	row := db.QueryRow(query, sentence)
	var jsonStr string
	if err := row.Scan(&jsonStr); err != nil {
		return nil, err
	}

	var tokens []string
	if err := json.Unmarshal([]byte(jsonStr), &tokens); err != nil {
		return nil, err
	}
	return tokens, nil
}

// Returns "easiest" sentence that contains word.
func PickSentence(db *sql.DB, word string) ([]string, error) {
	sentences, err := sentencesThatContain(db, word)
	if err != nil {
		return nil, err
	}

	ok := false
	sentence := 0
	minDifficulty := math.Inf(1)

	for id, difficulty := range sentences {
		if difficulty < minDifficulty {
			ok = true
			sentence = id
			minDifficulty = difficulty
		}
	}

	if !ok {
		panic("no sentences found")
	}
	return sentenceTokens(db, sentence)
}
