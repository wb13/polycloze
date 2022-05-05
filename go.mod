module github.com/lggruspe/polycloze-srs

go 1.18

replace github.com/mattn/go-sqlite3 => ../go-sqlite3

require (
	github.com/golang-migrate/migrate/v4 v4.15.2
	github.com/mattn/go-sqlite3 v1.14.13
)

require (
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	go.uber.org/atomic v1.7.0 // indirect
)
