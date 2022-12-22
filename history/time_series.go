// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package history

import (
	"time"
)

type Metric[T any] struct {
	Time  time.Time `json:"time"`
	Value T         `json:"value"`

	initialized bool
}

// Returns a series of zeros for given range and step size.
// Panics if the step size is < 1 second.
func Zeros[T any](from, to time.Time, step time.Duration) []Metric[T] {
	if step < time.Second {
		panic("only supports up to second precision")
	}

	var series []Metric[T]
	for current := from; current.Before(to); current = current.Add(step) {
		series = append(series, Metric[T]{
			Time: current,
		})
	}
	return series
}
