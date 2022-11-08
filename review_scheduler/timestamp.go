// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package review_scheduler

import (
	"time"
)

// Gets number of seconds in time.Duration as an int.
func seconds(d time.Duration) int {
	return int(d.Seconds())
}
