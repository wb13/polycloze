package flashcards

import (
	"database/sql"
	"math/rand"
	"strings"
	"sync"

	"github.com/lggruspe/polycloze/database"
	"github.com/lggruspe/polycloze/sentence_picker"
	"github.com/lggruspe/polycloze/translator"
	"github.com/lggruspe/polycloze/word_scheduler"
)

type Sentence struct {
	Id    int      // id in database
	Parts []string // Odd-numbered parts are blanks
}

type Item struct {
	Sentence    Sentence
	Translation string
}

type ItemGenerator struct {
	db *sql.DB
	tr translator.Translator

	// databases to be attached
	l1db          string
	l2db          string
	translationDb string
}

// Creates an ItemGenerator.
func NewItemGenerator(db *sql.DB, lang1Db, lang2Db, translationsDb string) (ItemGenerator, error) {
	var ig ItemGenerator
	ig.db = db

	// Initialize translator
	tr, err := translator.NewTranslator(db, lang2Db, lang1Db, translationsDb)
	if err != nil {
		return ig, err
	}
	ig.tr = *tr
	ig.l1db = lang1Db
	ig.l2db = lang2Db
	ig.translationDb = translationsDb
	return ig, nil
}

func getParts(tokens []string, word string) []string {
	var indices []int
	for i, token := range tokens {
		if strings.ToLower(token) == strings.ToLower(word) {
			indices = append(indices, i)
		}
	}

	if len(indices) == 0 {
		panic("something went wrong: Python casefold different from golang ToLower")
	}

	index := indices[rand.Intn(len(indices))]
	return []string{
		strings.Join(tokens[:index], ""),
		word,
		strings.Join(tokens[index+1:], ""),
	}
}

func (ig ItemGenerator) generateItem(word string) (Item, error) {
	var item Item

	session, err := database.NewSession(
		ig.db,
		ig.l1db,
		ig.l2db,
		ig.translationDb,
	)
	if err != nil {
		return item, err
	}

	sentence, err := sentence_picker.PickSentence(session, word)
	if err != nil {
		return item, err
	}
	translation, err := ig.tr.Translate(sentence.Text)
	if err != nil {
		return item, err
	}

	return Item{
		Translation: translation,
		Sentence: Sentence{
			Id:    sentence.Id,
			Parts: getParts(sentence.Tokens, word),
		},
	}, nil
}

// Generates up to n cloze items.
// Pass a negative value of n to get an unlimited number of items.
func (ig ItemGenerator) GenerateItems(n int) []Item {
	session, err := database.NewSession(
		ig.db,
		ig.l1db,
		ig.l2db,
		ig.translationDb,
	)
	if err != nil {
		return nil
	}

	words, err := word_scheduler.GetWords(session, n)
	if err != nil {
		return nil
	}

	session.Close()

	var wg sync.WaitGroup
	ch := make(chan Item, len(words))

	wg.Add(len(words))
	for _, word := range words {
		go func(ig *ItemGenerator, word string) {
			defer wg.Done()
			item, err := ig.generateItem(word)
			if err == nil {
				ch <- item
			}
		}(&ig, word)
	}

	wg.Wait()
	close(ch)

	var items []Item
	for item := range ch {
		items = append(items, item)
	}
	return items
}
