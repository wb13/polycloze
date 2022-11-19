// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package sentences

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/lggruspe/polycloze/database"
)

type Sentence struct {
	ID        int
	TatoebaID int64 // non-positive if none
	Text      string
	Tokens    []string
}

func findWordID[T database.Querier](q T, word string) (int, error) {
	query := `select id from word where word = ?`
	row := q.QueryRow(query, word)

	var id int
	err := row.Scan(&id)
	return id, err
}

func getSentence[T database.Querier](q T, id int) (*Sentence, error) {
	query := `select tatoeba_id, text, tokens from sentence where id = ?`
	row := q.QueryRow(query, id)

	var sentence Sentence
	sentence.ID = id
	var tatoebaID sql.NullInt64
	var jsonStr string

	err := row.Scan(&tatoebaID, &sentence.Text, &jsonStr)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(jsonStr), &sentence.Tokens); err != nil {
		return nil, err
	}

	if tatoebaID.Valid {
		sentence.TatoebaID = tatoebaID.Int64
	} else {
		sentence.TatoebaID = -1
	}
	return &sentence, nil
}

func PickSentence[T database.Querier](q T, word string) (*Sentence, error) {
	id, err := findWordID(q, word)
	if err != nil {
		return nil, err
	}

	// The course builder guarantees that all words have example sentences that
	// have the same difficulty (`frequency_class`) as the word.
	// Since the word scheduler only introduces words at the right difficulty,
	// the example sentences are also at the right difficulty.
	query := `
		SELECT id FROM contains
		JOIN sentence ON (sentence = id)
		WHERE word = ?
		ORDER BY random() LIMIT 1
	`
	row := q.QueryRow(query, id)
	var sentence int
	if err := row.Scan(&sentence); err != nil {
		return nil, err
	}
	return getSentence(q, sentence)
}

func Search[T database.Querier](q T, text string) (Sentence, error) {
	query := `
select id, tatoeba_id, tokens from sentence where text = ? collate nocase
`
	row := q.QueryRow(query, text)

	var sentence Sentence
	sentence.Text = text
	var tatoebaID sql.NullInt64
	var jsonStr string
	err := row.Scan(&sentence.ID, &tatoebaID, &jsonStr)
	if err != nil {
		return sentence, fmt.Errorf("sentence not found (%v): %v", text, err)
	}
	if tatoebaID.Valid {
		sentence.TatoebaID = tatoebaID.Int64
	} else {
		sentence.TatoebaID = -1
	}
	if err := json.Unmarshal([]byte(jsonStr), &sentence.Tokens); err != nil {
		return sentence, fmt.Errorf("sentence not found (%v): %v", text, err)
	}
	return sentence, nil
}
