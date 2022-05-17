package sentence_picker

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/lggruspe/polycloze/database"
)

var ErrNoSentenceFound error = errors.New("no sentence found")

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

func PickSentence(s *database.Session, word string, maxDifficulty int) (*Sentence, error) {
	id, err := findWordId(s, word)
	if err != nil {
		return nil, err
	}

	// Select sentence that contains word and isn't too "difficult"
	query := `
select id from l2.contains join l2.sentence on (contains.sentence = id)
where word = ? and frequency_class <= ?
order by random() limit 1
`
	row := s.QueryRow(query, id, maxDifficulty)

	var sentence int
	if err := row.Scan(&sentence); err != nil {
		return nil, err
	}
	return getSentence(s, sentence)
}

// Returns first sentence where the predicate value is true.
func FindSentence(s *database.Session, word string, maxDifficulty int, pred func(sentence *Sentence) bool) (*Sentence, error) {
	id, err := findWordId(s, word)
	if err != nil {
		return nil, err
	}

	// Select sentence that contains word and isn't too "difficult"
	query := `
select id from l2.contains join l2.sentence on (contains.sentence = id)
where word = ? and frequency_class <= ? order by random()
`
	rows, err := s.Query(query, id, maxDifficulty)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var sentence int
		if err := rows.Scan(&sentence); err != nil {
			continue
		}

		ps, err := getSentence(s, sentence)
		if err != nil {
			continue
		}

		if pred(ps) {
			return ps, nil
		}
	}
	return nil, ErrNoSentenceFound
}
