// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package flashcards

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/lggruspe/polycloze/text"
)

type Answer struct {
	// Text as it appears in the sentence.
	Text string `json:"text"`

	// Normalized form of the word
	Normalized string `json:"normalized"`
}

// Parts of a sentence.
// Represents a cloze if answers is non-empty.
type Part struct {
	Text    string   `json:"text"`
	Answers []Answer `json:"answers,omitempty"`
}

// Returns parts of cloze item.
func getParts(tokens []string, word string) []Part {
	word = text.Casefold(word)

	// Find all matching tokens.
	var indices []int
	for i, token := range tokens {
		if text.Casefold(token) == word {
			indices = append(indices, i)
		}
	}

	if len(indices) == 0 {
		message := fmt.Sprintf(
			"Python casefold different from golang casefold: %s, %v",
			word,
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
				Normalized: word,
			},
		},
	}
	return []Part{before, missing, after}
}
