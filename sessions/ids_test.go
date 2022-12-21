// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package sessions

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/lggruspe/polycloze/database"
)

func testDB() *sql.DB {
	db, err := database.OpenAuthDB(":memory:")
	if err != nil {
		panic(err)
	}
	return db
}

func TestReserveID(t *testing.T) {
	t.Parallel()
	id := "abcdefg"
	db := testDB()
	defer db.Close()

	if err := reserveID(db, id); err != nil {
		t.Fatal("expected ID to be available and err to be nil:", err)
	}

	err := reserveID(db, id)
	if err == nil {
		t.Fatal("expected uniqueness constraint error")
	}

	if !strings.Contains(err.Error(), "UNIQUE constraint failed") {
		t.Fatal("expected uniqueness constraint error")
	}
}

func TestDeleteIDNotExists(t *testing.T) {
	// It's not an error to delete an unused ID.
	t.Parallel()
	id := "abcdefg"
	db := testDB()
	defer db.Close()

	if err := deleteID(db, id); err != nil {
		t.Fatal("expected err to be nil:", err)
	}
}

func TestDeleteIDExists(t *testing.T) {
	// It should delete the ID from the DB.
	t.Parallel()
	id := "abcdefg"
	db := testDB()
	defer db.Close()

	if err := reserveID(db, id); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	if err := deleteID(db, id); err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	// The ID should be free to use again.
	if err := reserveID(db, id); err != nil {
		t.Fatal("expected err to be nil:", err)
	}
}

func BenchmarkGenerateUniqueID(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		db := testDB()
		defer db.Close()
		for pb.Next() {

			if _, err := generateUniqueID(db); err != nil {
				b.Log("expected err to be nil:", err)
			}
		}
	})
}
