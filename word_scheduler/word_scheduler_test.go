package word_scheduler

import (
	"testing"

	"github.com/lggruspe/polycloze/basedir"
	"github.com/lggruspe/polycloze/database"
	rs "github.com/lggruspe/polycloze/review_scheduler"
)

func init() {
	basedir.Init()
}

func wordScheduler() *database.Session {
	db, _ := rs.New(":memory:")
	s, err := database.NewSession(db, "", basedir.Language("spa"), "")
	if err != nil {
		panic(err)
	}
	return s
}

func TestFrequencyClass(t *testing.T) {
	s := wordScheduler()
	class := frequencyClass(s, "hola")
	if class <= 0 {
		t.Fatal("expected frequency class to be > 0")
	}
}

func TestCase(t *testing.T) {
	// Reviewed items should be auto-case-folded.
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
