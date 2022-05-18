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

func (ig ItemGenerator) GenerateItem(word string) (Item, error) {
	var item Item

	session, err := ig.Session()
	if err != nil {
		return item, err
	}
	defer session.Close()

	var translation string
	sentence, err := sentence_picker.FindSentence(session, word, 8, func(sent *sentence_picker.Sentence) bool {
		t, err := translator.Translate(session, sent.Text)
		if err != nil {
			return false
		}
		translation = t
		return true
	})
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

// Creates a cloze item for each word.
func (ig ItemGenerator) GenerateItems(words []string) []Item {
	var wg sync.WaitGroup
	ch := make(chan Item, len(words))

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
	close(ch)

	var items []Item
	for item := range ch {
		items = append(items, item)
	}
	return items
}
