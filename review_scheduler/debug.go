// Debugging stuff.
package review_scheduler

import (
	"fmt"
	"time"
)

func printReview(review *Review) {
	if review == nil {
		fmt.Println("Review(nil)")
		return
	}

	fmt.Printf(
		"Review(due=%v, interval=%v, reviewed=%v, correct=%v)\n",
		review.Due,
		review.Interval,
		review.Reviewed,
		review.Correct(),
	)
}

// Prints reviews in database for debugging purposes.
func printReviews[T CanQuery](db T) error {
	query := `
SELECT item, due, interval, reviewed FROM review`
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var item string
		var review Review
		var interval time.Duration

		var due string
		var reviewed string
		err := rows.Scan(
			&item,
			&due,
			&interval,
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
		review.Interval = interval * time.Second

		print(item, " ")
		printReview(&review)
	}
	return nil
}
