package translator

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/lggruspe/polycloze/database"
)

var ErrNoTranslationsFound = errors.New("no translations found")

type Sentence struct {
	Id        int
	TatoebaId int64 // negative if none
	Text      string
	Tokens    []string
}

func findSentence(s *database.Session, text string) (Sentence, error) {
	query := `
select id, tatoeba_id, tokens from sentence where text = ? collate nocase
`
	row := s.QueryRow(query, text)

	var sentence Sentence
	sentence.Text = text
	var tatoebaId sql.NullInt64
	var jsonStr string
	err := row.Scan(&sentence.Id, &tatoebaId, &jsonStr)
	if err != nil {
		return sentence, err
	}
	if tatoebaId.Valid {
		sentence.TatoebaId = tatoebaId.Int64
	} else {
		sentence.TatoebaId = -1
	}

	if err := json.Unmarshal([]byte(jsonStr), &sentence.Tokens); err != nil {
		return sentence, err
	}
	return sentence, nil
}

func Translate(s *database.Session, text string) (string, error) {
	sentence, err := findSentence(s, text)
	if err != nil || sentence.TatoebaId < 0 {
		return "", nil
	}
	query := `
select text from translation where tatoeba_id in
	(select target from translates where source = ?)
	order by random() limit 1
`
	row := s.QueryRow(query, sentence.TatoebaId)
	var translation string
	err = row.Scan(&translation)
	return translation, err
}
