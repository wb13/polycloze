// Database management stuff.
package review_scheduler

import (
	"database/sql"
	"strings"
	"time"
)

const layout string = "2006-01-02 15:04:05"

// Parses sqlite timestamps.
func parseTimestamp(timestamp string) (time.Time, error) {
	prefix := strings.TrimSpace(timestamp)[:len(layout)]
	return time.Parse(layout, prefix)
}

type CanQuery interface {
	*sql.DB | *sql.Tx

	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}
