package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/lggruspe/polycloze/database"
	"github.com/lggruspe/polycloze/translator"
)

func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	db, err := database.New(":memory:")
	if err != nil {
		log.Fatal(err)
	}

	session, err := database.NewSession(db, "../eng.db", "../spa.db", "../translations.db")
	if err != nil {
		log.Fatal(err)
	}

	text := readInput()
	result, err := translator.Translate(session, text)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}
