// Database management stuff.
package srs

import (
	"database/sql"
	"embed"
	"errors"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var fs embed.FS

const layout string = "2006-01-02 15:04:05"

// Upgrades the database to the latest version.
func migrateUp(db *sql.DB) error {
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

// Parses sqlite timestamps.
func parseTimestamp(timestamp string) (time.Time, error) {
	prefix := strings.TrimSpace(timestamp)[:len(layout)]
	return time.Parse(layout, prefix)
}
