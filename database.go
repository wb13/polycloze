package translator

import (
	"database/sql"
)

func attach(db *sql.DB, name, path string) error {
	query := `attach database ? as ?`
	_, err := db.Exec(query, path, name)
	return err
}
