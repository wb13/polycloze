package translator

import (
	"database/sql"
	"encoding/json"
	"errors"
	"math/rand"
)

var ErrNoTranslationsFound = errors.New("no translations found")

type Translator struct {
	db *sql.DB
}

type Sentence struct {
	Id        int
	TatoebaId int64 // negative if none
	Text      string
	Tokens    []string
}

// db should be empty in-memory database.
func NewTranslator(db *sql.DB, sourceDB, targetDB, translationDB string) (*Translator, error) {
	if err := attach(db, "source", sourceDB); err != nil {
		return nil, err
	}
	if err := attach(db, "target", targetDB); err != nil {
		return nil, err
	}
	if err := attach(db, "translation", translationDB); err != nil {
		return nil, err
	}
	return &Translator{db: db}, nil
}

func findSentence(db *sql.DB, text string) (Sentence, error) {
	query := `
select id, tatoeba_id, tokens from source.sentence where text = ?
collate nocase
`
	row := db.QueryRow(query, text)

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

// Returns tatoeba translations.
func (t Translator) tatoebaTranslate(text string) []string {
	sentence, err := findSentence(t.db, text)
	if err != nil || sentence.TatoebaId < 0 {
		return nil
	}
	query := `
select text from target.sentence where tatoeba_id in
	(select source from translation where target = ?
		union select target from translation where source = ?)
`
	rows, err := t.db.Query(query, sentence.TatoebaId, sentence.TatoebaId)
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

func (t Translator) Translate(sentence string) (string, error) {
	translations := t.tatoebaTranslate(sentence)
	// TODO use backup translation service if no translations found
	n := len(translations)
	if n == 0 {
		return "", ErrNoTranslationsFound
	}
	choice := rand.Intn(n)
	return translations[choice], nil
}
