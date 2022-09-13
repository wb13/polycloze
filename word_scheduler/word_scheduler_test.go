// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package word_scheduler

import (
	"database/sql"
	"testing"

	rs "github.com/lggruspe/polycloze/review_scheduler"
	"github.com/lggruspe/polycloze/utils"
)

func wordScheduler() *sql.DB {
	return utils.TestingDatabase()
}

func TestFrequencyClass(t *testing.T) {
	t.Parallel()

	s := wordScheduler()

	query := `insert into word (word, frequency_class) values (?, ?)`
	if _, err := s.Exec(query, "foo", 1); err != nil {
		panic(err)
	}
	if _, err := s.Exec(query, "bar", 2); err != nil {
		panic(err)
	}

	if class := frequencyClass(s, "Foo"); class != 1 {
		t.Fatal("expected frequency class to be 1")
		// should be case-insensitive
	}

	if class := frequencyClass(s, "bar"); class != 2 {
		t.Fatal("expected frequency class to be 2")
	}

	if class := frequencyClass(s, "baz"); class != 0 {
		t.Fatal("expected frequency class to be 0")
		// if not in database
	}
}

func TestCase(t *testing.T) {
	// Reviewed items should be auto-case-folded.
	t.Parallel()

	s := wordScheduler()

	if err := UpdateWord(s, "Foo", false); err != nil {
		t.Fatal("expected err to be nil", err)
	}

	words, err := rs.ScheduleReviewNow(s, 1)
	if err != nil {
		t.Fatal("expected err to be nil", err)
	}

	if len(words) != 1 {
		t.Fatal("expected words to contain one item")
	}

	if words[0] != "foo" {
		t.Error("expected word to be \"foo\"")
	}
}
