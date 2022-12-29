// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package replay

import (
	"encoding/csv"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type ReviewEvent struct {
	Word     string
	Reviewed time.Time
	Correct  bool
}

// Turns the event into a CSV record.
func (e ReviewEvent) Record() []string {
	reviewed := strconv.FormatInt(e.Reviewed.Unix(), 10)
	correct := "1"
	if !e.Correct {
		correct = "0"
	}
	return []string{e.Word, reviewed, correct}
}

type ReviewReader struct {
	csvReader *csv.Reader
}

func NewReviewReader(r *csv.Reader) *ReviewReader {
	return &ReviewReader{csvReader: r}
}

func (r *ReviewReader) ReadReview() (ReviewEvent, error) {
	record, err := r.csvReader.Read()
	if err != nil {
		return ReviewEvent{}, fmt.Errorf("failed to read review from CSV: %v", err)
	}
	if len(record) != 3 {
		return ReviewEvent{}, errors.New(
			"failed to read review from CSV: incorrect number of fields",
		)
	}

	i, err := strconv.ParseInt(record[1], 10, 64)
	if err != nil {
		return ReviewEvent{}, fmt.Errorf("failed to read review from CSV: %v", err)
	}

	var correct bool
	switch record[2] {
	case "0":
		correct = false
	case "1":
		correct = true
	default:
		return ReviewEvent{}, errors.New(
			"failed to read review from CSV: invalid correct value",
		)
	}

	return ReviewEvent{
		Word:     record[0],
		Reviewed: time.Unix(i, 0),
		Correct:  correct,
	}, nil
}

type ReviewWriter struct {
	csvWriter *csv.Writer
}

func NewReviewWriter(r *csv.Writer) *ReviewWriter {
	return &ReviewWriter{csvWriter: r}
}

func (w *ReviewWriter) WriteReview(e ReviewEvent) error {
	if err := w.csvWriter.Write(e.Record()); err != nil {
		return fmt.Errorf("failed to write review into CSV: %v", err)
	}
	w.csvWriter.Flush()
	return nil
}
