package sentence_picker

import (
	"database/sql"
	"encoding/json"

	"github.com/lggruspe/polycloze/database"
)

type Sentence struct {
	Id        int
	TatoebaId int64
	Text      string
	Tokens    []string
}

func findWordId(s *database.Session, word string) (int, error) {
	query := `select id from word where word = ?`
	row := s.QueryRow(query, word)

	var id int
	err := row.Scan(&id)
	return id, err
}

func getSentence(s *database.Session, id int) (*Sentence, error) {
	query := `select tatoeba_id, text, tokens from sentence where id = ?`
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
	// I.e. sentence.frequency_class <= student.frequency_class, or if there's no
	// such sentence, returns the sentence with the minimum frequency_class instead.
	query := `
select coalesce(
	(select id from contains join sentence on (sentence = id)
		where word = ? and frequency_class <= ?
		order by random() limit 1),
	(select coalesce(id, min(frequency_class)) from contains join sentence on (sentence = id)
		where word = ?))
`
	row := s.QueryRow(query, id, maxDifficulty, id)

	var sentence int
	if err := row.Scan(&sentence); err != nil {
		return nil, err
	}
	return getSentence(s, sentence)
}
