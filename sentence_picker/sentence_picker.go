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

// Returns "easiest" sentence that contains word.
func PickSentence(s *database.Session, word string) (*Sentence, error) {
	id, err := findWordId(s, word)
	if err != nil {
		return nil, err
	}

	// Sentence difficulty = max word difficulty
	// NOTE Only takes easiest sentence from random sample of sentences.
	// Taking the easiest sentence over all the sentences may take too long.
	query := `
select sentence, min(difficulty) from
	(select sentence, max(difficulty) as difficulty from l2.contains
		join word_difficulty using (word)
		where sentence in
			(select sentence from l2.contains where word = ? order by random() limit 500)
		group by sentence)
`
	var sentence int
	var difficulty float64
	row := s.QueryRow(query, id)
	if err := row.Scan(&sentence, &difficulty); err != nil {
		return nil, err
	}

	return getSentence(s, sentence)
}
