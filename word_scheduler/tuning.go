// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

// More tuning stuff, depends on course tables.
package word_scheduler

import (
	"github.com/lggruspe/polycloze/database"
)

// Smallest frequency class of unseen word.
func easiestUnseen(s *database.Session) int {
	query := `
select min(frequency_class) from word
where word not in (select item from review)
`

	var difficulty int
	row := s.QueryRow(query)
	_ = row.Scan(&difficulty)
	return difficulty
}

// Make sure student level is not lower than lowest unseen word.
func postTune(s *database.Session) error {
	easiest := easiestUnseen(s)
	preferred := PreferredDifficulty(s)

	if preferred < easiest {
		query := `
update student set
	frequency_class = ?,
	correct = 0,
	incorrect = 0
`
		_, err := s.Exec(query, easiest)
		return err
	}
	return nil
}
