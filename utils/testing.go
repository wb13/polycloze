// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

// Testing utils.
package utils

import (
	"database/sql"
	"embed"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/polycloze/polycloze/database"
)

//go:embed test.sql
var fs embed.FS

func isBeginTransaction(query string) bool {
	return strings.HasPrefix(query, "begin transaction") || strings.HasPrefix(query, "BEGIN TRANSACTION")
}

func isCommit(query string) bool {
	return strings.HasPrefix(query, "commit") || strings.HasPrefix(query, "COMMIT")
}

// Returns list of statements in test.sql.
func readTestSQL() []string {
	bytes, err := fs.ReadFile("test.sql")
	if err != nil {
		panic(err)
	}

	var queries []string

	substrings := strings.Split(string(bytes), ";")
	for _, substring := range substrings {
		query := strings.TrimSpace(substring)
		if strings.HasPrefix(query, "--") {
			continue
		}
		if isBeginTransaction(query) {
			continue
		}
		if isCommit(query) {
			continue
		}
		queries = append(queries, query)
	}
	return queries
}

// Returns DB for testing.
// NOTE Caller has to Close the db.
func TestingDatabase() *sql.DB {
	db, err := database.OpenReviewDB(":memory:")
	if err != nil {
		panic(err)
	}

	for _, query := range readTestSQL() {
		if _, err := db.Exec(query); err != nil {
			panic(err)
		}
	}
	return db
}
