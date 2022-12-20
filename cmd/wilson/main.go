// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/lggruspe/polycloze/wilson"
)

func experiment(ns, nf int, z float64) {
	lower := wilson.Wilson(ns, nf, z)
	upper := wilson.Wilson(ns, nf, -z)
	fmt.Printf("%.4f <= p\n", lower)
	fmt.Printf("p <= %.4f\n", upper)
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
	confidence := []float64{0.8, 0.9, 0.95, 0.99, 0.999}
	zs := []float64{-0.845, -1.285, -1.645, -2.325, -3.1}
	// z scores are for one-sided lower bounds.
	// Negate for one-sided upper bound z-scores.

	for i, z := range zs {
		a := confidence[i]
		fmt.Println("confidence:", a)
		experiment(ns, nf, z)
	}
}
