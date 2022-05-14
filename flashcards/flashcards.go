package flashcards

import (
	"database/sql"
	"math/rand"
	"strings"

	"github.com/lggruspe/polycloze-sentence-picker/sentence_picker"
	"github.com/lggruspe/polycloze-srs/word_scheduler"
	"github.com/lggruspe/polycloze-translator/translator"
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
}

// Creates an ItemGenerator.
//
// reviewDb doesn't have to exist, but the other databases are expected to
// exist and to have the appropriate schema.
func NewItemGenerator(reviewDb, lang1Db, lang2Db, translationsDb string) (ItemGenerator, error) {
	var ig ItemGenerator

	// Initialize word_scheduler
	db, err := word_scheduler.New(reviewDb, lang2Db)
	if err != nil {
		return ig, err
	}
	ig.db = db

	// Initialize sentence_picker
	err = sentence_picker.InitSentencePicker(db, lang2Db, reviewDb)
	if err != nil {
		return ig, err
	}

	// Initialize translator
	tr, err := translator.NewTranslator(db, lang2Db, lang1Db, translationsDb)
	if err != nil {
		return ig, err
	}
	ig.tr = *tr
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

	sentence, err := sentence_picker.PickSentence(ig.db, word)
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
	words, err := word_scheduler.GetWords(ig.db, n)
	if err != nil {
		return nil
	}

	var items []Item
	for _, word := range words {
		item, err := ig.generateItem(word)
		if err == nil {
			items = append(items, item)
		}
	}
	return items
}
