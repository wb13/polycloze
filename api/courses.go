// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package api

import (
	"fmt"
	"path/filepath"

	"github.com/lggruspe/polycloze/basedir"
)

type Course struct {
	l1, l2 Language
}

func AvailableCourses() []Course {
	var courses []Course

	glob := "[a-z][a-z][a-z]-[a-z][a-z][a-z].db"
	matches, _ := filepath.Glob(filepath.Join(basedir.DataDir, glob))
	for _, match := range matches {
		course, err := getCourseInfo(match)
		if err == nil {
			courses = append(courses, course)
		}
	}
	return courses
}

// Input: path to course db file.
func getCourseInfo(path string) (Course, error) {
	var course Course

	db, err := openDB(path)
	if err != nil {
		return course, err
	}
	defer db.Close()

	query := `select id, code, name from language`
	rows, err := db.Query(query)
	if err != nil {
		return course, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, code, name string
		if err := rows.Scan(&id, &code, &name); err != nil {
			return course, err
		}

		switch id {
		case "l1":
			course.l1.Code = code
			course.l1.Name = name
		case "l2":
			course.l2.Code = code
			course.l2.Name = name
		}
	}

	if course.l1.Code == "" || course.l2.Code == "" {
		return course, fmt.Errorf("invalid course database: %s\n", path)
	}
	return course, nil
}
