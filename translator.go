package translator

import (
	"database/sql"
)

type Translator struct {
	db *sql.DB
	// TODO backup translation service
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

func (t Translator) Translate(sentence string) (string, error) {
	// TODO result should vary if there are multiple translations
	// TODO use backup translation service if no translations found
	query := `
select target.sentence.text from translation
join source.sentence on (translation.source = source.sentence.id)
join target.sentence on (translation.target = target.sentence.id)
where source.sentence.text = ?
`
	row := t.db.QueryRow(query, sentence)
	var translation string
	err := row.Scan(&translation)
	return translation, err
}
