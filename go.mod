module github.com/lggruspe/polycloze-flashcards

go 1.18

replace github.com/mattn/go-sqlite3 => ../go-sqlite3

replace github.com/lggruspe/polycloze-srs => ../polycloze-srs

replace github.com/lggruspe/polycloze-sentence-picker => ../polycloze-sentence-picker

replace github.com/lggruspe/polycloze-translator => ../polycloze-translator

require github.com/mattn/go-sqlite3 v1.14.13 // indirect

require (
	github.com/lggruspe/polycloze-sentence-picker v0.0.0
	github.com/lggruspe/polycloze-srs v0.0.0
	github.com/lggruspe/polycloze-translator v0.0.0
)

require (
	github.com/golang-migrate/migrate/v4 v4.15.2 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	go.uber.org/atomic v1.7.0 // indirect
)
