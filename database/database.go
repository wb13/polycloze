// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

// Database management stuff.
package database

import (
	"context"
	"database/sql"
	"embed"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

//go:embed migrations
var migrations embed.FS

func init() {
	goose.SetBaseFS(migrations)
	goose.SetLogger(goose.NopLogger())
	if err := goose.SetDialect("sqlite3"); err != nil {
		panic(err)
	}
}

// NOTE Caller has to Close the db.
func Open(path string) (*sql.DB, error) {
	return sql.Open("sqlite3", path)
}

// Attaches database to the connection.
func attach(con *sql.Conn, name, path string) error {
	query := `attach database ? as ?`
	_, err := con.ExecContext(context.TODO(), query, path, name)
	return err
}

// Detaches database from connection.
func detach(con *sql.Conn, name string) error {
	query := `detach database ?`
	_, err := con.ExecContext(context.TODO(), query, name)
	return err
}

type Querier interface {
	*sql.DB | *sql.Tx | *Connection

	Begin() (*sql.Tx, error)
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}
