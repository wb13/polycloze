// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package flashcards

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"sync"

	"github.com/lggruspe/polycloze/database"
	"github.com/lggruspe/polycloze/sentences"
	"github.com/lggruspe/polycloze/text"
	"github.com/lggruspe/polycloze/translator"
	"github.com/lggruspe/polycloze/word_scheduler"
)

// Different from sentences.Sentence
type Sentence struct {
	ID        int      `json:"id"`    // id in database
	Parts     []string `json:"parts"` // Odd-numbered parts are blanks
	TatoebaID int64    `json:"tatoebaID,omitempty"`
}

type Item struct {
	Sentence    Sentence               `json:"sentence"`
	Translation translator.Translation `json:"translation"`
}

type ItemGenerator struct {
	db       *sql.DB
	courseDB string // to be attached
}

// NOTE Caller has to close connection.
func (ig ItemGenerator) Session() (*database.Session, error) {
	return database.NewSession(ig.db, ig.courseDB)
}

// Creates an ItemGenerator.
func NewItemGenerator(db *sql.DB, courseDB string) ItemGenerator {
	return ItemGenerator{
		db:       db,
		courseDB: courseDB,
	}
}

func getParts(tokens []string, word string) []string {
	var indices []int
	for i, token := range tokens {
		if text.Casefold(token) == text.Casefold(word) {
			indices = append(indices, i)
		}
	}

	if len(indices) == 0 {
		message := fmt.Sprintf("Python casefold different from golang x case folder: %s, %v", word, tokens)
		panic(message)
	}

	index := indices[rand.Intn(len(indices))]
	return []string{
		strings.Join(tokens[:index], ""),
		tokens[index],
		strings.Join(tokens[index+1:], ""),
	}
}

func generateWords(
	db *sql.DB,
	n int,
	pred func(word string) bool,
	hooks ...database.ConnectionHook,
) ([]string, error) {
	ctx := context.TODO()
	con, err := database.NewConnection(db, ctx)
	if err != nil {
		return nil, err
	}
	defer con.Close()

	for _, hook := range hooks {
		if err := hook.Enter(con); err != nil {
			return nil, err
		}
		defer func(hook database.ConnectionHook) {
			_ = hook.Exit(con)
		}(hook)
	}
	return word_scheduler.GetWordsWith(con, n, pred)
}

func GenerateItem[T database.Querier](q T, word string) (Item, error) {
	var item Item

	sentence, err := sentences.PickSentence(q, word, word_scheduler.PreferredDifficulty(q))
	if err != nil {
		return item, err
	}

	translation, err := translator.Translate(q, *sentence)
	if err != nil {
		panic(fmt.Errorf("could not translate sentence (%v): %v", sentence, err))
	}
	return Item{
		Translation: translation,
		Sentence: Sentence{
			ID:        sentence.ID,
			Parts:     getParts(sentence.Tokens, word),
			TatoebaID: sentence.TatoebaID,
		},
	}, nil
}

// Creates a cloze item for each word.
func GenerateItems(db *sql.DB, words []string, hooks ...database.ConnectionHook) []Item {
	ch := make(chan Item, len(words))
	GenerateItemsIntoChannel(db, ch, words, hooks...)
	close(ch)

	items := make([]Item, 0)
	for item := range ch {
		items = append(items, item)
	}
	return items
}

// Enter functions are invoked in a FIFO manner, while Exit functions are deferred.
func GenerateItemsIntoChannel(
	db *sql.DB,
	ch chan Item,
	words []string,
	hooks ...database.ConnectionHook,
) {
	var wg sync.WaitGroup
	wg.Add(len(words))

	// TODO use request context instead
	ctx := context.TODO()

	for _, word := range words {
		go func(word string) {
			defer wg.Done()

			con, err := database.NewConnection(db, ctx)
			if err != nil {
				return
			}
			defer con.Close()

			for _, hook := range hooks {
				if err := hook.Enter(con); err != nil {
					return
				}
				defer func(hook database.ConnectionHook) {
					_ = hook.Exit(con)
				}(hook)
			}

			if item, err := GenerateItem(con, word); err == nil {
				ch <- item
			}
		}(word)
	}
	wg.Wait()
}

// Returns list of flashcards to show.
// n: max number of flashcards to return.
func Get(
	db *sql.DB,
	n int,
	pred func(word string) bool,
	hooks ...database.ConnectionHook,
) []Item {
	words, err := generateWords(db, n, pred, hooks...)
	if err != nil {
		return nil
	}
	return GenerateItems(db, words, hooks...)
}
