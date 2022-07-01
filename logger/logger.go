package logger

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lggruspe/polycloze/basedir"
)

func timestamp() string {
	layout := "2006-01-02 15:04:05"
	return time.Now().UTC().Format(layout)
}

func prefix(correct bool) string {
	result := "x"
	if correct {
		result = "/"
	}
	return fmt.Sprintf("%v %v ", result, timestamp())
}

func LogReview(l2 string, correct bool, word string) error {
	logFile := basedir.Log(l2)
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	logger := log.New(f, prefix(correct), 0)
	logger.Println(word)
	return nil
}
