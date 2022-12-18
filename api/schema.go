// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

// Schema definitions of JSON requests sent to and from the server.
package api

import (
	"github.com/lggruspe/polycloze/difficulty"
	"github.com/lggruspe/polycloze/flashcards"
)

type ReviewResult struct {
	Word    string `json:"word"`
	Correct bool   `json:"correct"`
}

// JSON request schema.
type FlashcardsRequest struct {
	Limit      int                   `json:"limit"`
	Difficulty difficulty.Difficulty `json:"difficulty"`
	Reviews    []ReviewResult        `json:"reviews"`
}

// JSON response schema.
type FlashcardsResponse struct {
	Items      []flashcards.Item     `json:"items"`
	Difficulty difficulty.Difficulty `json:"difficulty"`
}
