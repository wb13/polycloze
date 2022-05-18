// Buffered item generator.
package buffer

import (
	"github.com/lggruspe/polycloze/flashcards"
)

type BufferedItem struct {
	word string
	item flashcards.Item
}

type ItemBuffer struct {
	Channel    chan BufferedItem
	Words      map[string]bool
	BufferSize int

	ig flashcards.ItemGenerator
}

func NewItemBuffer(ig flashcards.ItemGenerator) ItemBuffer {
	size := 150
	return ItemBuffer{
		Channel:    make(chan BufferedItem, size),
		Words:      make(map[string]bool),
		BufferSize: size,
		ig:         ig,
	}
}

func (buf *ItemBuffer) Add(x BufferedItem) {
	buf.Channel <- x
	buf.Words[x.word] = true
}

// Fetch items and store in buffer.
func (buf *ItemBuffer) Fetch() error {
	n := cap(buf.Channel) - len(buf.Channel)
	words, err := buf.ig.GenerateWordsWith(n, func(word string) bool {
		_, ok := buf.Words[word]
		return !ok
	})
	if err != nil {
		return err
	}

	for _, word := range words {
		// TODO use goroutines?
		item, err := buf.ig.GenerateItem(word)
		if err == nil {
			buf.Add(BufferedItem{word: word, item: item})
		}
	}
	return nil
}

// Take an item out of buffer.
func (buf *ItemBuffer) Take() flashcards.Item {
	bufItem := <-buf.Channel
	delete(buf.Words, bufItem.word)
	return bufItem.item
}
