// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package translator

import (
	"errors"
	"fmt"

	"github.com/lggruspe/polycloze/database"
	"github.com/lggruspe/polycloze/sentences"
)

type Translation struct {
	TatoebaID int64  `json:"tatoebaID,omitempty"` // non-positive if none
	Text      string `json:"text"`
}

func Translate[T database.Querier](q T, sentence sentences.Sentence) (Translation, error) {
	var translation Translation

	if sentence.TatoebaID <= 0 {
		return translation, errors.New("sentence has no TatoebaID")
	}

	query := `
		SELECT tatoeba_id, text FROM translation
		WHERE tatoeba_id = (
			SELECT target FROM translates
			WHERE source = ?
			ORDER BY random() LIMIT 1
		)
	`

	row := q.QueryRow(query, sentence.TatoebaID)
	err := row.Scan(&translation.TatoebaID, &translation.Text)
	if err != nil {
		return translation, fmt.Errorf("failed to translate sentence: %v", err)
	}
	return translation, err
}
