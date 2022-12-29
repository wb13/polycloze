// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package replay

import (
	"testing"
	"time"
)

func TestParseLineCorrect(t *testing.T) {
	t.Parallel()

	event, err := ParseLine("/ 2020-01-01 00:00:00 test")
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	if !event.Correct {
		t.Fatal("expected event.Correct to be true")
	}
}

func TestParseLineIncorrect(t *testing.T) {
	t.Parallel()

	event, err := ParseLine("x 2020-01-01 00:00:00 test")
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	if event.Correct {
		t.Fatal("expected event.Correct to be false")
	}
}

func TestParseLineWord(t *testing.T) {
	t.Parallel()

	event, err := ParseLine("/ 2020-01-01 00:00:00 Foo bar")
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	if event.Word != "Foo bar" {
		t.Fatal("expected event.Word to be equal to 'Foo bar':", event.Word)
	}
}

func TestParse(t *testing.T) {
	t.Parallel()

	log := `/ 2021-01-01 01:02:03 foo
x 2022-02-02 01:02:03 bar
`
	events, err := Parse(log)
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	expected := []LogEvent{
		{
			Correct:   true,
			Timestamp: time.Date(2021, 0o1, 0o1, 1, 2, 3, 0, time.UTC),
			Word:      "foo",
		},
		{
			Correct:   false,
			Timestamp: time.Date(2022, 0o2, 0o2, 1, 2, 3, 0, time.UTC),
			Word:      "bar",
		},
	}

	if len(events) != len(expected) {
		t.Fatalf("expected events to be %v: %v\n", expected, events)
	}

	for i, event := range events {
		if event != expected[i] {
			t.Fatalf("expected events to be %v: %v\n", expected, events)
		}
	}
}

func TestParseWithComment(t *testing.T) {
	t.Parallel()

	log := `/ 2022-01-01 00:00:00 foo
# ignore me
x 2022-01-02 00:00:00 foo
`
	events, err := Parse(log)
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	if len(events) != 2 {
		t.Fatal("expected events to contain two elements:", events)
	}
}

func TestParseWithBlank(t *testing.T) {
	t.Parallel()

	log := `/ 2022-01-01 00:00:00 foo

x 2022-01-02 00:00:00 foo
`
	events, err := Parse(log)
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	if len(events) != 2 {
		t.Fatal("expected events to contain two elements:", events)
	}
}

func TestString(t *testing.T) {
	t.Parallel()

	line := `/ 2020-01-01 00:00:00 test`

	event, err := ParseLine(line)
	if err != nil {
		t.Fatal("expected err to be nil:", err)
	}

	output := event.String()
	if line != output {
		t.Fatal(
			"expected ParseLine and LogEvent.String to be inverse functions",
			line,
			output,
		)
	}
}
