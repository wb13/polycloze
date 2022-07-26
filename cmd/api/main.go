// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/lggruspe/polycloze/api"
)

func main() {
	config := api.Config{AllowCORS: true}
	r, err := api.Router(config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
