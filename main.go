package main

import (
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/lggruspe/polycloze-srs/review_scheduler"
)

func assertNil(value any) {
	if value != nil {
		log.Fatal(value)
	}
}

func main() {
	db, err := review_scheduler.New("test.db")
	assertNil(err)

	assertNil(review_scheduler.UpdateReview(db, "foo", false))
	assertNil(review_scheduler.UpdateReview(db, "foo", true))
	assertNil(review_scheduler.UpdateReview(db, "bar", true))

	items, err := review_scheduler.ScheduleReviewNow(db, -1)
	assertNil(err)
	for _, item := range items {
		fmt.Println(item)
	}

	/*
	database.UpgradeFile("review.db", "migrations/review_scheduler")

	db, _ := sql.Open("sqlite3", ":memory:")
	database.Attach(db, "review", "review.db")
	database.Attach(db, "language", "language.db")
	*/
}
