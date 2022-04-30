// Debugging stuff.
package srs

import (
	"database/sql"
)

// Prints Reviews in database for debugging purposes.
func printReviews(tx *sql.Tx) error {
	query := `
SELECT word, due, interval, reviewed, correct, streak FROM Review`
	rows, err := tx.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var word string
		var review Review
		err := rows.Scan(
			&word,
			&review.Due,
			&review.Interval,
			&review.Reviewed,
			&review.Correct,
			&review.Streak,
		)
		if err != nil {
			return err
		}

		println(
			"Review(word={}, due={}, interval={}, reviewed={}, correct={}, streak={})",
			word,
			review.Due.String(),
			review.Interval,
			review.Reviewed.String(),
			review.Correct,
			review.Streak,
		)
	}
	return nil
}
