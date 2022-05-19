package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/lggruspe/polycloze/api"
)

func main() {
	config := api.Config{
		ReviewDb:      "review.db",
		Lang1Db:       "../eng.db",
		Lang2Db:       "../spa.db",
		TranslationDb: "../translations.db",
		AllowCORS:     true,
	}

	mux, err := api.Mux(config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", mux))
}
