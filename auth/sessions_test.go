// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package auth

import (
	"testing"
)

func TestReserveID(t *testing.T) {
	t.Parallel()
	db := openDB()
	defer db.Close()

	if err := reserveID(db, "abc"); err != nil {
		t.Fatal("expected ID to be available:", err)
	}
	if err := reserveID(db, "abc"); err == nil {
		t.Fatal("expected ID to be taken")
	}
}
