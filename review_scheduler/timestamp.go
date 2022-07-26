// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package review_scheduler

import (
	"strings"
	"time"
)

const layout string = "2006-01-02 15:04:05"

func formatTime(t time.Time) string {
	return t.Format(layout)
}

// Parses sqlite timestamps.
func parseTimestamp(timestamp string) (time.Time, error) {
	prefix := strings.TrimSpace(timestamp)[:len(layout)]
	return time.Parse(layout, prefix)
}

// Gets number of seconds in time.Duration as an int.
func seconds(d time.Duration) int {
	return int(d.Seconds())
}
