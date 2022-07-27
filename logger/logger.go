// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

const layout string = "2006-01-02 15:04:05"

func timestamp() string {
	return time.Now().UTC().Format(layout)
}

func prefix(correct bool) string {
	result := "x"
	if correct {
		result = "/"
	}
	return fmt.Sprintf("%v %v ", result, timestamp())
}

func LogReview(file string, correct bool, word string) error {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}

	logger := log.New(f, prefix(correct), 0)
	logger.Println(word)
	return nil
}
