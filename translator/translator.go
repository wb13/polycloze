// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package translator

import (
	"errors"

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
select tatoeba_id, text from translation where tatoeba_id in
	(select target from translates where source = ?)
	order by random() limit 1
`
	row := q.QueryRow(query, sentence.TatoebaID)
	err := row.Scan(&translation.TatoebaID, &translation.Text)
	return translation, err
}
