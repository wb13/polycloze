package api

import (
	"github.com/lggruspe/polycloze/basedir"
)

type LanguageStats struct {
	// all-time
	Seen  int `json:"seen"`
	Total int `json:"total"`

	// today
	Learned  int `json:"learned"`
	Reviewed int `json:"reviewed"`

	// today
	Correct   int `json:"correct"`
	Incorrect int `json:"incorrect"`
}

func queryInt(path, query string) (int, error) {
	var result int

	db, err := openDb(path)
	if err != nil {
		return result, err
	}
	defer db.Close()

	row := db.QueryRow(query)
	err = row.Scan(&result)
	return result, err
}

func countSeen(lang string) (int, error) {
	return queryInt(basedir.Review(lang), `select count(*) from review`)
}

// Total count of words in lang (given as three-letter code).
func countTotal(lang string) (int, error) {
	return queryInt(basedir.Language(lang), `select count(*) from word`)
}

// New words learned today.
func countLearnedToday(lang string) (int, error) {
	query := `select count(*) from review where learned >= current_date`
	return queryInt(basedir.Review(lang), query)
}

// Number of words reviewed today, excluding new words.
func countReviewedToday(lang string) (int, error) {
	query := `
select count(*) from review where reviewed >= current_date
and learned < current_date
`
	return queryInt(basedir.Review(lang), query)
}

// Number of correct answers today.
func countCorrectToday(lang string) (int, error) {
	query := `select count(*) from review where reviewed >= current_date and correct`
	return queryInt(basedir.Review(lang), query)
}

// Number of incorrect answers today.
func countIncorrectToday(lang string) (int, error) {
	query := `select count(*) from review where reviewed >= current_date and not(correct)`
	return queryInt(basedir.Review(lang), query)
}

func getLanguageStats(lang string) (*LanguageStats, error) {
	seen, err := countSeen(lang)
	if err != nil {
		return nil, err
	}

	total, err := countTotal(lang)
	if err != nil {
		return nil, err
	}

	learned, err := countLearnedToday(lang)
	if err != nil {
		return nil, err
	}

	reviewed, err := countReviewedToday(lang)
	if err != nil {
		return nil, err
	}

	correct, err := countCorrectToday(lang)
	if err != nil {
		return nil, err
	}

	incorrect, err := countIncorrectToday(lang)
	if err != nil {
		return nil, err
	}

	return &LanguageStats{
		Seen:      seen,
		Total:     total,
		Learned:   learned,
		Reviewed:  reviewed,
		Correct:   correct,
		Incorrect: incorrect,
	}, nil
}
