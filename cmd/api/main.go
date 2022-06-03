package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/lggruspe/polycloze/api"
)

func main() {
	config := api.Config{
		L1:        "eng",
		L2:        "spa",
		AllowCORS: true,
	}

	mux, err := api.Mux(config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", mux))
}
