package main

import (
	"database/sql"
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
	db, err := sql.Open("sqlite3", "test.db")
	assertNil(err)

	rs, err := review_scheduler.InitReviewScheduler(db)
	assertNil(err)

	assertNil(rs.Update("foo", false))
	assertNil(rs.Update("foo", true))
	assertNil(rs.Update("bar", true))

	items, err := rs.ScheduleNow(-1)
	assertNil(err)
	for _, item := range items {
		fmt.Println(item)
	}
}
