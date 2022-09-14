// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package database

import (
	"context"
	"database/sql"
	"testing"
)

func database() *sql.DB {
	db, err := Open(":memory:")
	if err != nil {
		panic(err)
	}
	return db
}

func connection(db *sql.DB) *Connection {
	con, err := NewConnection(db, context.TODO())
	if err != nil {
		panic(err)
	}
	return con
}

func TestDefaultConnectionHook(t *testing.T) {
	// Default Enter and Exit functions shouldn't simply return nil.
	hook := DefaultConnectionHook()

	db := database()
	defer db.Close()

	con := connection(db)
	defer con.Close()

	if err := hook.Enter(con); err != nil {
		t.Fatal("expected err to be nil:", err)
	}
	if err := hook.Exit(con); err != nil {
		t.Fatal("expected err to be nil:", err)
	}
}
