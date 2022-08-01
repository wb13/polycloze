// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package translator

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/lggruspe/polycloze/database"
)

var ErrNoTranslationsFound = errors.New("no translations found")

type Sentence struct {
	ID        int
	TatoebaID int64 // non-positive if none
	Text      string
	Tokens    []string
}

type Translation struct {
	TatoebaID int64  `json:"tatoebaID,omitempty"` // non-positive if none
	Text      string `json:"text"`
}

func findSentence(s *database.Session, text string) (Sentence, error) {
	query := `
select id, tatoeba_id, tokens from sentence where text = ? collate nocase
`
	row := s.QueryRow(query, text)

	var sentence Sentence
	sentence.Text = text
	var tatoebaID sql.NullInt64
	var jsonStr string
	err := row.Scan(&sentence.ID, &tatoebaID, &jsonStr)
	if err != nil {
		return sentence, err
	}
	if tatoebaID.Valid {
		sentence.TatoebaID = tatoebaID.Int64
	} else {
		sentence.TatoebaID = -1
	}

	if err := json.Unmarshal([]byte(jsonStr), &sentence.Tokens); err != nil {
		return sentence, err
	}
	return sentence, nil
}

func Translate(s *database.Session, text string) (Translation, error) {
	var translation Translation

	sentence, err := findSentence(s, text)
	if err != nil || sentence.TatoebaID < 0 {
		return translation, fmt.Errorf("sentence not found: %v", text)
	}
	query := `
select tatoeba_id, text from translation where tatoeba_id in
	(select target from translates where source = ?)
	order by random() limit 1
`
	row := s.QueryRow(query, sentence.TatoebaID)
	err = row.Scan(&translation.TatoebaID, &translation.Text)
	return translation, err
}
