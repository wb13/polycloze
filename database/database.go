// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

// Database management stuff.
package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var fs embed.FS

func init() {
	goose.SetBaseFS(fs)
	if err := goose.SetDialect("sqlite3"); err != nil {
		panic(err)
	}
}

// NOTE Caller has to Close the db.
func Open(path string) (*sql.DB, error) {
	return sql.Open("sqlite3", path)
}

// Convenience function for creating upgraded sqlite DB.
func New(path string) (*sql.DB, error) {
	db, err := Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	if err := Upgrade(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to upgrade database: %v", err)
	}
	return db, nil
}

// Upgrades database to the latest version.
func Upgrade(db *sql.DB) error {
	return goose.Up(db, "migrations")
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
