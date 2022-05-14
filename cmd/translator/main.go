package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	translator "github.com/lggruspe/polycloze/translator"
)

func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}

	tr, err := translator.NewTranslator(db, "spa.db", "eng.db", "translations.db")
	if err != nil {
		log.Fatal(err)
	}

	text := readInput()
	result, err := tr.Translate(text)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}
