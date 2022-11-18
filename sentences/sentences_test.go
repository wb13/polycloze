// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package sentences

import (
	"context"
	"testing"

	"github.com/lggruspe/polycloze/basedir"
	"github.com/lggruspe/polycloze/database"
)

func BenchmarkPickSentence(b *testing.B) {
	db, err := database.New(":memory:")
	if err != nil {
		b.Fatal("expected err to be nil:", err)
	}
	defer db.Close()

	con, err := database.NewConnection(
		db,
		context.Background(),
		database.AttachCourse(basedir.Course("eng", "deu")),
	)
	if err != nil {
		b.Fatal("expected err to be nil:", err)
	}
	defer con.Close()

	for i := 0; i < b.N; i++ {
		sentence, err := PickSentence(con, "was", 20)
		if err != nil {
			b.Fatal("expected err to be nil:", err)
		}
		b.Log("result:", sentence)
	}
}
