package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/lggruspe/polycloze/basedir"
	"github.com/lggruspe/polycloze/database"
	"github.com/lggruspe/polycloze/translator"
)

func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func main() {
	if err := basedir.Init(); err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().UnixNano())

	db, err := database.New(":memory:")
	if err != nil {
		log.Fatal(err)
	}

	session, err := database.NewSession(
		db,
		basedir.Language("eng"),
		basedir.Language("spa"),
		basedir.Translation("eng", "spa"),
	)
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
