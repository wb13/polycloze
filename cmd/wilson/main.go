// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	rs "github.com/lggruspe/polycloze/review_scheduler"
)

func experiment(ns, nf int, z float64) {
	lower := rs.Wilson(ns, nf, z)
	upper := rs.Wilson(ns, nf, -z)
	fmt.Printf("%.4f <= p <= 1\n", lower)
	fmt.Printf("0 <= p <= %.4f\n", upper)
	fmt.Println()
}

func atoi(s string) int {
	x, _ := strconv.Atoi(s)
	return x
}

func parseArgs() (int, int) {
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		log.Fatal("incomplete arguments: ns, nf")
	}
	return atoi(flag.Arg(0)), atoi(flag.Arg(1))
}

func main() {
	ns, nf := parseArgs()
	confidence := []float64{0.8, 0.9, 0.95, 0.99}
	zs := []float64{-0.845, -1.285, -1.645, -2.325}

	for i, z := range zs {
		a := confidence[i]
		fmt.Println("confidence:", a)
		experiment(ns, nf, z)
	}
}
