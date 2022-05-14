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

// Returns map of sentence ID -> sentence difficulty of sentences containing word.
func sentencesThatContain(s *database.Session, word string) (map[int]float64, error) {
	// NOTE seen counter is only used to add variation in the returned sentences
	query := `
select contains.sentence as sentence,
			count(word) + sum(difficulty) + coalesce(counter, 0.0) as difficulty
from l2.contains
join word_difficulty using(word)
left join seen on (contains.sentence = seen.sentence)
where contains.sentence in
	(select sentence from l2.contains join l2.word on (contains.word = word.id) where word.word = ?)
group by contains.sentence;
`
	rows, err := s.Query(query, word)
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

func IncrementSeenCount(db *sql.DB, sentence int) error {
	query := `
update or ignore seen
set last = current_timestamp, counter = counter + 1
where sentence = ?
`
	_, err := db.Exec(query, sentence)
	return err
}
