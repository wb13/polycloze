package sentence_picker

import (
	"database/sql"
	"encoding/json"
	"math"

	"github.com/lggruspe/polycloze/database"
)

type Sentence struct {
	Id        int
	TatoebaId int64
	Text      string
	Tokens    []string
}

func findWordId(s *database.Session, word string) (int, error) {
	query := `select id from l2.word where word = ?`
	row := s.QueryRow(query, word)

	var id int
	err := row.Scan(&id)
	return id, err
}

// Returns map of sentence ID -> sentence difficulty of sentences containing word.
func sentencesThatContain(s *database.Session, word string) (map[int]float64, error) {
	id, err := findWordId(s, word)
	if err != nil {
		return nil, err
	}

	// TODO won't the same sentence just get chosen repeatedly?
	query := `
select contains.sentence as sentence, max(difficulty) as difficulty
from l2.contains join word_difficulty using(word)
where contains.sentence in (select sentence from l2.contains where word = ?)
group by contains.sentence;
`
	rows, err := s.Query(query, id)
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

func getSentence(s *database.Session, id int) (*Sentence, error) {
	query := `select tatoeba_id, text, tokens from l2.sentence where id = ?`
	row := s.QueryRow(query, id)

	var sentence Sentence
	sentence.Id = id
	var tatoebaId sql.NullInt64
	var jsonStr string

	err := row.Scan(&tatoebaId, &sentence.Text, &jsonStr)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(jsonStr), &sentence.Tokens); err != nil {
		return nil, err
	}

	if tatoebaId.Valid {
		sentence.TatoebaId = tatoebaId.Int64
	} else {
		sentence.TatoebaId = -1
	}
	return &sentence, nil
}

// Returns "easiest" sentence that contains word.
func PickSentence(s *database.Session, word string) (*Sentence, error) {
	sentences, err := sentencesThatContain(s, word)
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
	return getSentence(s, sentence)
}
