// Database management stuff.
package database

import (
	"database/sql"
	"embed"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*/*.sql
var fs embed.FS

// Upgrades database to the latest version.
func Upgrade(db *sql.DB, migrationsPath string) error {
	dbDriver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return err
	}

	srcDriver, err := iofs.New(fs, migrationsPath)
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

// Attaches database.
func Attach(db *sql.DB, name, path string) error {
	query := `attach database ? as ?`
	_, err := db.Exec(query, path, name)
	return err
}
