// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package database

import (
	"context"
	"database/sql"
	"errors"
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
	t.Parallel()

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

func TestConnectionHookEnterExit(t *testing.T) {
	// Enter and Exit hooks should be called in the proper order.
	t.Parallel()

	var result string

	appendA := ConnectionHook{
		Enter: func(c *Connection) error {
			result += "A"
			return nil
		},
		Exit: func(c *Connection) error {
			result += "a"
			return nil
		},
	}

	appendB := ConnectionHook{
		Enter: func(c *Connection) error {
			result += "B"
			return nil
		},
		Exit: func(c *Connection) error {
			result += "b"
			return nil
		},
	}

	db := database()
	defer db.Close()

	con, err := NewConnection(db, context.TODO(), appendA, appendB)
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	if result != "AB" {
		t.Fatal("expected appendA.Enter and appendB.Enter to have been called")
	}

	if err := con.Close(); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	if result != "ABba" {
		t.Fatal("expected appendB.Exit and appendA.Exit to have been called")
	}
}

func TestNewConnectionError(t *testing.T) {
	// Exit hooks should be called immediately in the proper order on failure.
	t.Parallel()

	enter := 0
	exit := 0

	count := ConnectionHook{
		Enter: func(c *Connection) error {
			enter++
			return nil
		},
		Exit: func(c *Connection) error {
			exit++
			return nil
		},
	}

	fail := ConnectionHook{
		Enter: func(c *Connection) error {
			return errors.New(":(")
		},
		Exit: func(c *Connection) error {
			return errors.New(":(")
		},
	}

	db := database()
	defer db.Close()

	con, err := NewConnection(db, context.TODO(), count, count, fail)
	if err == nil {
		t.Fatal("expected err to be not nil")
	}
	if con != nil {
		con.Close()
		t.Fatal("expected con to be nil")
	}

	if enter <= 0 {
		t.Fatal("expected at least some of the Enter hooks to be called")
	}
	if exit <= 0 {
		t.Fatal("expected at least some of the Exit hooks to be called")
	}
	if enter != exit {
		t.Fatal("expected the same number of Enter and Exit to be called")
	}
}
