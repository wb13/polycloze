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
	prerunChecks()

	args := parseArgs()
	config := api.Config{AllowCORS: args.cors, Port: args.port}
	r, err := api.Router(config)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listening on port %v\n", args.port)
	log.Printf("Start learning: http://127.0.0.1:%v\n", args.port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", args.port), r))
}

func prerunChecks() {
	courses := api.AvailableCourses()
	if len(courses) == 0 {
		log.Fatal("Couldn't find installed courses. Please visit https://github.com/lggruspe/polycloze-data")
	}
}
