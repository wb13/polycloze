// Review "sessions."
package database

import (
	"context"
	"database/sql"
)

// Returns a connection with the necessary attached databases.
//
// NOTE Caller is expected to close the connection after use.
func NewSession(db *sql.DB, l1db, l2db, translationDb string) (*sql.Conn, error) {
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
	return con, nil
}
