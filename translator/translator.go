// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package translator

import (
	"fmt"

	"github.com/lggruspe/polycloze/database"
	"github.com/lggruspe/polycloze/sentences"
)

type Translation struct {
	TatoebaID int64  `json:"tatoebaID,omitempty"` // non-positive if none
	Text      string `json:"text"`
}

func Translate(s *database.Session, sentence sentences.Sentence) (Translation, error) {
	var translation Translation

	if sentence.TatoebaID <= 0 {
		return translation, fmt.Errorf("could not translate sentence (%v)", sentence)
	}

	query := `
select tatoeba_id, text from translation where tatoeba_id in
	(select target from translates where source = ?)
	order by random() limit 1
`
	row := s.QueryRow(query, sentence.TatoebaID)
	if err := row.Scan(&translation.TatoebaID, &translation.Text); err != nil {
		return translation, fmt.Errorf("could not translate sentence (%v): %v", sentence, err.Error())
	}
	return translation, nil
}
