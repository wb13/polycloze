// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/lggruspe/polycloze/api"
)

type Args struct {
	cors bool
	port int
}

func parseArgs() Args {
	var args Args

	flag.BoolVar(&args.cors, "c", false, "allow CORS")
	flag.IntVar(&args.port, "p", 3000, "port number")
	flag.Parse()
	return args
}

func main() {
	args := parseArgs()

	config := api.Config{AllowCORS: args.cors, Port: args.port}
	r, err := api.Router(config)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listening on port %v\n", args.port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", args.port), r))
}
