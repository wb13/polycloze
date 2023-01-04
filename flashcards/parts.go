// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package flashcards

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/polycloze/polycloze/text"
	"github.com/polycloze/polycloze/word_scheduler"
)

type Answer struct {
	// Text as it appears in the sentence.
	Text string `json:"text"`

	// Normalized form of the word
	Normalized string `json:"normalized"`

	// Is the word new/previously unseen?
	New bool `json:"new"`

	// Only has to be meaningful for new words.
	Difficulty int `json:"difficulty"`
}

// Parts of a sentence.
// Represents a cloze if answers is non-empty.
type Part struct {
	Text    string   `json:"text"`
	Answers []Answer `json:"answers,omitempty"`
}

// Returns parts of cloze item.
func getParts(tokens []string, word word_scheduler.Word) []Part {
	// TODO word: string -> Word
	normalized := text.Casefold(word.Word)

	// Find all matching tokens.
	var indices []int
	for i, token := range tokens {
		if text.Casefold(token) == normalized {
			indices = append(indices, i)
		}
	}

	if len(indices) == 0 {
		message := fmt.Sprintf(
			"Python casefold different from golang casefold: %s, %v",
			normalized,
			tokens,
		)
		panic(message)
	}

	// Pick a random one if there are multiple matches.
	// TODO turn all matching tokens into blanks instead.
	index := indices[rand.Intn(len(indices))]

	before := Part{
		Text: strings.Join(tokens[:index], ""),
	}
	after := Part{
		Text: strings.Join(tokens[index+1:], ""),
	}

	missing := Part{
		Text: tokens[index],
		Answers: []Answer{
			{
				Text:       tokens[index],
				Normalized: normalized,
				New:        word.New,
				Difficulty: word.Difficulty,
			},
		},
	}
	return []Part{before, missing, after}
}
