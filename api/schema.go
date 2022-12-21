// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

// Schema definitions of JSON requests sent to and from the server.
package api

import (
	"github.com/lggruspe/polycloze/difficulty"
	"github.com/lggruspe/polycloze/flashcards"
	"github.com/lggruspe/polycloze/review_scheduler"
)

type ReviewResult = review_scheduler.Result

// JSON request schema.
type FlashcardsRequest struct {
	Limit      int                    `json:"limit"`
	Difficulty *difficulty.Difficulty `json:"difficulty"`
	Reviews    []ReviewResult         `json:"reviews"`
	Exclude    []string               `json:"exclude"`

	// Sometimes used by client if for some reason they can't pass the token via
	// HTTP headers (e.g. `sendBeacon`).
	CSRFToken string `json:"csrfToken"`
}

// JSON response schema.
type FlashcardsResponse struct {
	Items      []flashcards.Item      `json:"items"`
	Difficulty *difficulty.Difficulty `json:"difficulty"`
}

type SetCourseRequest struct {
	L1Code string `json:"l1Code"`
	L2Code string `json:"l2Code"`
}

type SetCourseResponse struct {
	Ok bool `json:"ok"`
}
