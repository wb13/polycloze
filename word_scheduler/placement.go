// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package word_scheduler

import (
	"github.com/lggruspe/polycloze/database"
	"github.com/lggruspe/polycloze/wilson"
)

// Estimates student's level with a frequency class.
func Placement[T database.Querier](q T) int {
	var class int

	query := `
		SELECT frequency_class, correct, incorrect
		FROM new_word_stat
		ORDER BY frequency_class ASC
	`

	rows, err := q.Query(query)
	if err != nil {
		return class
	}
	defer rows.Close()

	var correct, incorrect int
	for rows.Next() {
		var _class int
		if err := rows.Scan(&_class, &correct, &incorrect); err != nil {
			goto done
		}

		if wilson.IsTooHard(correct, incorrect) {
			goto done
		}

		class = _class
		if !wilson.IsTooEasy(correct, incorrect) {
			goto done
		}
	}

	// This should only be reachable if all rows were visited.
	if wilson.IsTooEasy(correct, incorrect) {
		class += 1
	}

done:
	// Estimated level shouldn't be lower than the lowest unseen frequency class.
	// Otherwise, the estimated level will be stuck if there's a frequency class
	// with a low enough score.
	if easiest := easiestUnseen(q); easiest > class {
		class = easiest
	}

	return class
}

// Smallest frequency class of unseen word.
func easiestUnseen[T database.Querier](q T) int {
	query := `
select min(frequency_class) from word
where word not in (select item from review)
`

	var difficulty int
	row := q.QueryRow(query)
	_ = row.Scan(&difficulty)
	return difficulty
}

// This should be called only when an item is seen for the first time.
func updateNewWordStat[T database.Querier](q T, frequencyClass int, correct bool) error {
	var x, y int
	if correct {
		x = 1
	} else {
		y = 1
	}

	query := `
		INSERT INTO new_word_stat (frequency_class, correct, incorrect)
		VALUES (?, ?, ?)
		ON CONFLICT (frequency_class) DO UPDATE SET
			correct = correct + excluded.correct,
			incorrect = incorrect + excluded.incorrect
	`
	_, err := q.Exec(query, frequencyClass, x, y)
	return err
}
