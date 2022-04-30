// Debugging stuff.
package srs

import (
	"fmt"
)

func printReview(review Review) {
	fmt.Printf(
		"Review(due=%v, interval=%v, reviewed=%v, correct=%v, streak=%v)\n",
		review.Due,
		review.Interval,
		review.Reviewed,
		review.Correct,
		review.Streak,
	)
}

// Prints Reviews in database for debugging purposes.
func printReviews[T CanQuery](db T) error {
	query := `
SELECT id, word, due, interval, reviewed, correct, streak FROM Review`
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var word string
		var review Review

		var due string
		var reviewed string
		err := rows.Scan(
			&id,
			&word,
			&due,
			&review.Interval,
			&reviewed,
			&review.Correct,
			&review.Streak,
		)
		if err != nil {
			return err
		}

		parsedDue, err := parseTimestamp(due)
		if err != nil {
			return err
		}
		parsedReviewed, err := parseTimestamp(reviewed)
		if err != nil {
			return err
		}

		review.Due = parsedDue
		review.Reviewed = parsedReviewed

		print(id, " ")
		printReview(review)
	}
	return nil
}
