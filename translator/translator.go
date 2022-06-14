package translator

import (
	"database/sql"
	"encoding/json"
	"errors"
	"math/rand"

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
select id, tatoeba_id, tokens from l2.sentence where text = ?
collate nocase
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

func IsReversed(s *database.Session) (bool, error) {
	var l1, l2 string

	query := `select code from l1.info`
	row := s.QueryRow(query)
	if err := row.Scan(&l1); err != nil {
		return false, err
	}

	query = `select code from l2.info`
	row = s.QueryRow(query)
	if err := row.Scan(&l2); err != nil {
		return false, err
	}

	if l1 == l2 {
		panic("invalid language pair")
	}

	return l2 < l1, nil
}

// Returns tatoeba translations.
func tatoebaTranslate(s *database.Session, text string) []string {
	sentence, err := findSentence(s, text)
	if err != nil || sentence.TatoebaId < 0 {
		return nil
	}
	query := `
select text from l1.sentence where tatoeba_id in
	(select l1 from translation where l2 = ?)
`
	if reversed, _ := IsReversed(s); reversed {
		query = `
select text from l1.sentence where tatoeba_id in
	(select l2 from translation where l1 = ?)
`
	}

	rows, err := s.Query(query, sentence.TatoebaId)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var translations []string
	for rows.Next() {
		var translation string
		if err := rows.Scan(&translation); err == nil {
			translations = append(translations, translation)
		}
	}
	return translations
}

func Translate(s *database.Session, sentence string) (string, error) {
	translations := tatoebaTranslate(s, sentence)
	// TODO use backup translation service if no translations found
	n := len(translations)
	if n == 0 {
		return "", ErrNoTranslationsFound
	}
	choice := rand.Intn(n)
	return translations[choice], nil
}
