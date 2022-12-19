// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package review_scheduler

// Review results
type Result struct {
	Word    string `json:"word"`
	Correct bool   `json:"correct"`
}
