package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/lggruspe/polycloze/api"
)

func parseArgs() api.Config {
	if len(os.Args) < 2 {
		log.Fatal("missing L1")
	}
	if len(os.Args) < 3 {
		log.Fatal("missing L2")
	}

	l1 := os.Args[1]
	l2 := os.Args[2]

	ok1 := false
	ok2 := false
	for _, lang := range api.SupportedLanguages() {
		if ok1 && ok2 {
			break
		}
		if !ok1 && l1 == lang.Code {
			ok1 = true
		}
		if !ok2 && l2 == lang.Code {
			ok2 = true
		}
	}

	return api.Config{
		L1:        l1,
		L2:        l2,
		AllowCORS: true,
	}
}

func main() {
	config := parseArgs()
	r, err := api.Router(config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
