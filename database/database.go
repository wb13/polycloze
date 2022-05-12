// Database management stuff.
package database

import (
	"database/sql"
	"embed"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*/*.sql
var fs embed.FS

// Upgrades database (specified by path) to the latest version.
func Upgrade(databaseUrl string, migrationsPath string) error {
	srcDriver, err := iofs.New(fs, migrationsPath)
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance(
		"iofs",
		srcDriver,
		databaseUrl,
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
