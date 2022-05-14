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

	query := `drop view if exists word_difficulty`
	if _, err := s.Exec(query); err != nil {
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

	session := Session{con: con}
	query := `
create temp view word_difficulty as
select frequency_class.id as word,
			 frequency_class/(1.0 + coalesce(level, 0.0)) as difficulty
from l2.frequency_class left join most_recent_review on (frequency_class.word = most_recent_review.item)
`
	if _, err := session.Exec(query); err != nil {
		return nil, err
	}
	return &session, nil
}
