// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package translator

import (
	"errors"
	"fmt"

	"github.com/lggruspe/polycloze/database"
	"github.com/lggruspe/polycloze/sentences"
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

func Translate(s *database.Session, text string) (Translation, error) {
	var translation Translation

	sentence, err := sentences.Search(s, text)
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
