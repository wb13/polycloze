// Database management stuff.
package database

import (
	"database/sql"
	"embed"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*.sql
var fs embed.FS

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

// Upgrades database (specified by path) to the latest version.
func UpgradeFile(dbPath string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	return Upgrade(db)
}

// Attaches database.
func Attach(db *sql.DB, name, path string) error {
	query := `attach database ? as ?`
	_, err := db.Exec(query, path, name)
	return err
}
