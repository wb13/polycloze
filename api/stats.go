// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package api

import (
	"fmt"

	"github.com/lggruspe/polycloze/basedir"
	"github.com/lggruspe/polycloze/database"
)

// If upgrade is non-empty, upgrades the database.
func queryInt(path, query string, upgrade ...bool) (int, error) {
	var result int

	db, err := database.Open(path)
	if err != nil {
		return 0, fmt.Errorf("could not open db (%v) for query (%v): %v", path, query, err)
	}
	defer db.Close()

	if len(upgrade) > 0 {
		if err := database.Upgrade(db); err != nil {
			return result, err
		}
	}

	row := db.QueryRow(query)
	err = row.Scan(&result)
	return result, err
}

// Total count of words in course.
func CountTotal(l1, l2 string) (int, error) {
	return queryInt(basedir.Course(l1, l2), `select count(*) from word`)
}
