// Review "sessions."
package database

import (
	"context"
	"database/sql"
)

type Session struct {
	con *sql.Conn
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
func NewSession(db *sql.DB, l1db, l2db, translationDb string) (*Session, error) {
	ctx := context.TODO()
	con, err := db.Conn(ctx)
	if err != nil {
		return nil, err
	}

	if err := attach(con, "l1", l1db); err != nil {
		return nil, err
	}
	if err := attach(con, "l2", l2db); err != nil {
		return nil, err
	}
	if err := attach(con, "translation", translationDb); err != nil {
		return nil, err
	}
	return &Session{con: con}, nil
}
