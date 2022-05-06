// Debugging stuff.
package srs

import (
	"fmt"
)

func printReview(review *Review) {
	if review == nil {
		fmt.Println("Review(nil)")
		return
	}

	fmt.Printf(
		"Review(due=%v, interval=%v, reviewed=%v, correct=%v, level=%v)\n",
		review.Due,
		review.Interval,
		review.Reviewed,
		review.Correct(),
		review.Level(),
	)
}

// Prints Reviews in database for debugging purposes.
func printReviews[T CanQuery](db T) error {
	query := `
SELECT id, item, due, interval, reviewed FROM Review`
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var item string
		var review Review

		var due string
		var reviewed string
		err := rows.Scan(
			&id,
			&item,
			&due,
			&review.Interval,
			&reviewed,
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

		print(id, " ", item, " ")
		printReview(&review)
	}
	return nil
}
