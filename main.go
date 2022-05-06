package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/lggruspe/polycloze-srs/srs"
)

func assertNil(value any) {
	if value != nil {
		log.Fatal(value)
	}
}

func main() {
	db, err := sql.Open("sqlite3", "test.db")
	assertNil(err)

	ws, err := srs.InitWordScheduler(db)
	assertNil(err)

	assertNil(ws.Update("foo", false))
	assertNil(ws.Update("foo", true))
	assertNil(ws.Update("bar", true))

	items, err := ws.ScheduleNow(-1)
	assertNil(err)
	for _, item := range items {
		fmt.Println(item)
	}
}
