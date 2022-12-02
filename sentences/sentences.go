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
	ID int `json:"id,omitempty"`

	// Non-positive if none
	TatoebaID int64    `json:"tatoebaID"`
	Text      string   `json:"text"`
	Tokens    []string `json:"tokens,omitempty"`
}

func findWordID[T database.Querier](q T, word string) (int, error) {
	query := `select id from word where word = ?`
	row := q.QueryRow(query, word)

	var id int
	err := row.Scan(&id)
	return id, err
}

func PickSentence[T database.Querier](q T, word string) (Sentence, error) {
	id, err := findWordID(q, word)
	if err != nil {
		return Sentence{}, err
	}

	// The course builder guarantees that all words have example sentences that
	// have the same difficulty (`frequency_class`) as the word.
	// Since the word scheduler only introduces words at the right difficulty,
	// the example sentences are also at the right difficulty.
	query := `
		SELECT id, tatoeba_id, text, tokens FROM contains
		JOIN sentence ON (sentence = id)
		WHERE word = ?
		ORDER BY random() LIMIT 1
	`
	row := q.QueryRow(query, id)

	var sentence Sentence
	var tatoebaID sql.NullInt64
	var tokens string

	err = row.Scan(&sentence.ID, &tatoebaID, &sentence.Text, &tokens)
	if err != nil {
		return sentence, err
	}

	if err := json.Unmarshal([]byte(tokens), &sentence.Tokens); err != nil {
		return sentence, err
	}

	if tatoebaID.Valid {
		sentence.TatoebaID = tatoebaID.Int64
	} else {
		sentence.TatoebaID = -1
	}
	return sentence, nil
}

// Returns random sentence from the database.
// The results don't include tokens.
// NOTE Only picks random sentence from first 10,000 sentences in the DB for
// speed.
func RandomSentences[T database.Querier](q T, limit int) ([]Sentence, error) {
	query := `
		SELECT id, tatoeba_id, text
		FROM (SELECT * FROM sentence LIMIT 10000)
		ORDER BY random()
		LIMIT ?
	`
	rows, err := q.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to pick random sentences: %v", err)
	}
	defer rows.Close()

	var sentences []Sentence
	for rows.Next() {
		var sentence Sentence
		var tatoebaID sql.NullInt64

		if err := rows.Scan(&sentence.ID, &tatoebaID, &sentence.Text); err != nil {
			return nil, fmt.Errorf("failed to pick random sentences: %v", err)
		}

		if tatoebaID.Valid {
			sentence.TatoebaID = tatoebaID.Int64
		} else {
			sentence.TatoebaID = -1
		}
		sentences = append(sentences, sentence)
	}
	return sentences, nil
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
