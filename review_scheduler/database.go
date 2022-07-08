// Database management stuff.
package review_scheduler

import (
	"database/sql"
)

type CanQuery interface {
	*sql.DB | *sql.Tx

	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}
