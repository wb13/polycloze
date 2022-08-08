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
	TatoebaID int64
	Text      string
	Tokens    []string
}

func findWordID(s *database.Session, word string) (int, error) {
	query := `select id from word where word = ?`
	row := s.QueryRow(query, word)

	var id int
	err := row.Scan(&id)
	return id, err
}

func getSentence(s *database.Session, id int) (*Sentence, error) {
	query := `select tatoeba_id, text, tokens from sentence where id = ?`
	row := s.QueryRow(query, id)

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

func PickSentence(s *database.Session, word string, maxDifficulty int) (*Sentence, error) {
	id, err := findWordID(s, word)
	if err != nil {
		return nil, err
	}

	// Select sentence that contains word and isn't too "difficult"
	// I.e. sentence.frequency_class <= student.frequency_class, or if there's no
	// such sentence, returns the sentence with the minimum frequency_class instead.
	query := `
select coalesce(
	(select id from contains join sentence on (sentence = id)
		where word = ? and frequency_class <= ?
		order by random() limit 1),
	(select coalesce(id, min(frequency_class)) from contains join sentence on (sentence = id)
		where word = ?))
`
	row := s.QueryRow(query, id, maxDifficulty, id)

	var sentence int
	if err := row.Scan(&sentence); err != nil {
		return nil, err
	}
	return getSentence(s, sentence)
}

func Search(s *database.Session, text string) (Sentence, error) {
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
		return sentence, fmt.Errorf("sentence not found (%v): %v", text, err.Error())
	}
	if tatoebaID.Valid {
		sentence.TatoebaID = tatoebaID.Int64
	} else {
		sentence.TatoebaID = -1
	}
	if err := json.Unmarshal([]byte(jsonStr), &sentence.Tokens); err != nil {
		return sentence, fmt.Errorf("sentence not found (%v): %v", text, err.Error())
	}
	return sentence, nil
}
