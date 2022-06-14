package flashcards

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"sync"

	"golang.org/x/text/cases"

	"github.com/lggruspe/polycloze/database"
	"github.com/lggruspe/polycloze/sentence_picker"
	"github.com/lggruspe/polycloze/translator"
	"github.com/lggruspe/polycloze/word_scheduler"
)

type Sentence struct {
	Id    int      `json:"id"`    // id in database
	Parts []string `json:"parts"` // Odd-numbered parts are blanks
}

type Item struct {
	Sentence    Sentence `json:"sentence"`
	Translation string   `json:"translation"`
}

type ItemGenerator struct {
	db *sql.DB

	// databases to be attached
	l1db          string
	l2db          string
	translationDb string
}

func (ig ItemGenerator) Session() (*database.Session, error) {
	return database.NewSession(
		ig.db,
		ig.l1db,
		ig.l2db,
		ig.translationDb,
	)
}

// Creates an ItemGenerator.
func NewItemGenerator(db *sql.DB, lang1Db, lang2Db, translationDb string) ItemGenerator {
	return ItemGenerator{
		db:            db,
		l1db:          lang1Db,
		l2db:          lang2Db,
		translationDb: translationDb,
	}
}

func getParts(tokens []string, word string) []string {
	var indices []int
	for i, token := range tokens {
		caser := cases.Fold()
		if caser.String(token) == caser.String(word) {
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

func (ig ItemGenerator) GenerateItem(word string) (Item, error) {
	var item Item

	session, err := ig.Session()
	if err != nil {
		return item, err
	}
	defer session.Close()

	sentence, err := sentence_picker.FindTranslatedSentence(session, word, 8)
	if err != nil {
		return item, err
	}

	translation, err := translator.Translate(session, sentence.Text)
	if err != nil {
		panic("unexpected error")
	}
	return Item{
		Translation: translation,
		Sentence: Sentence{
			Id:    sentence.Id,
			Parts: getParts(sentence.Tokens, word),
		},
	}, nil
}

// Generates up to n words.
// Pass a negative value of n to get an unlimited number of items.
func (ig ItemGenerator) GenerateWords(n int) ([]string, error) {
	session, err := ig.Session()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	return word_scheduler.GetWords(session, n)
}

// Same as GenerateWords, but takes an additional predicate argument.
// Only returns words that satisfy the predicate.
func (ig ItemGenerator) GenerateWordsWith(n int, pred func(word string) bool) ([]string, error) {
	session, err := ig.Session()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	return word_scheduler.GetWordsWith(session, n, pred)
}

// Creates a cloze item for each word.
func (ig ItemGenerator) GenerateItems(words []string) []Item {
	ch := make(chan Item, len(words))
	ig.GenerateItemsIntoChannel(ch, words)
	close(ch)

	var items []Item
	for item := range ch {
		items = append(items, item)
	}
	return items
}

func (ig ItemGenerator) GenerateItemsIntoChannel(ch chan Item, words []string) {
	var wg sync.WaitGroup
	wg.Add(len(words))
	for _, word := range words {
		go func(ig *ItemGenerator, word string) {
			defer wg.Done()
			item, err := ig.GenerateItem(word)
			if err == nil {
				ch <- item
			}
		}(&ig, word)
	}
	wg.Wait()
}
