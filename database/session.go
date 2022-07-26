// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

// Review "sessions."
package database

import (
	"context"
	"database/sql"
)

type Session struct {
	con *sql.Conn
}

func (s *Session) Exec(query string, args ...any) (sql.Result, error) {
	return s.con.ExecContext(context.TODO(), query, args...)
}

func (s *Session) Query(query string, args ...any) (*sql.Rows, error) {
	return s.con.QueryContext(context.TODO(), query, args...)
}

func (s *Session) QueryRow(query string, args ...any) *sql.Row {
	return s.con.QueryRowContext(context.TODO(), query, args...)
}

func (s *Session) Begin() (*sql.Tx, error) {
	return s.con.BeginTx(context.TODO(), nil)
}

func (s *Session) Close() error {
	if err := detach(s.con, "translation"); err != nil {
		return err
	}
	if err := detach(s.con, "l1"); err != nil {
		return err
	}
	if err := detach(s.con, "l2"); err != nil {
		return err
	}
	return s.con.Close()
}

// Returns a connection with the necessary attached databases.
//
// NOTE Caller is expected to close the connection after use.
func NewSession(db *sql.DB, courseDB string) (*Session, error) {
	ctx := context.TODO()
	con, err := db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	if err := attach(con, "course", courseDB); err != nil {
		return nil, err
	}
	return &Session{con: con}, nil
}
