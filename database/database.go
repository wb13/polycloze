// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

// Database management stuff.
package database

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*.sql
var fs embed.FS

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
	dbDriver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return err
	}

	srcDriver, err := iofs.New(fs, "migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance(
		"iofs",
		srcDriver,
		"sqlite3",
		dbDriver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
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
	*sql.DB | *sql.Tx | *Connection | *Session

	Begin() (*sql.Tx, error)
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}
