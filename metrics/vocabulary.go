// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package metrics

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type VocabularyData struct {
	data     map[int][]int
	nSamples int
	start    time.Time
	end      time.Time
}

func (d *VocabularyData) add(interval, i, total int) {
	if d.data[interval] == nil {
		d.data[interval] = make([]int, d.nSamples)
	}
	d.data[interval][i] = total
}

type vocabularyDatasetSchema struct {
	Name string `json:"name"`
	Data []int  `json:"data"`
}

type vocabularyDataSchema struct {
	Start    time.Time                 `json:"start"`
	End      time.Time                 `json:"end"`
	NSamples int                       `json:"nSamples"`
	Datasets []vocabularyDatasetSchema `json:"datasets"`
}

// Encodes vocabulary data to JSON.
func (d VocabularyData) EncodeJSON() ([]byte, error) {
	var datasets []vocabularyDatasetSchema
	for key, values := range d.data {
		dataset := vocabularyDatasetSchema{
			Name: fmt.Sprintf("%vh", key),
			Data: values,
		}
		datasets = append(datasets, dataset)
	}
	return json.Marshal(vocabularyDataSchema{
		Start:    d.start,
		End:      d.end,
		NSamples: d.nSamples,
		Datasets: datasets,
	})
}

func newVocabularyData(start, end time.Time, nSamples int) VocabularyData {
	return VocabularyData{
		data:     make(map[int][]int),
		nSamples: nSamples,
		start:    start,
		end:      end,
	}
}

func CountVocabulary(db *sql.DB, start, end time.Time, nSamples int) (VocabularyData, error) {
	query := `
		SELECT interval, total, CAST(@n * (t - @start) / CAST(@end - @start AS REAL) AS INTEGER)
		FROM interval_count
		WHERE @start <= t AND t < @end
		ORDER BY t
	`
	rows, err := db.Query(
		query,
		sql.Named("n", nSamples),
		sql.Named("start", start.Unix()),
		sql.Named("end", end.Unix()),
	)
	if err != nil {
		return VocabularyData{}, err
	}
	defer rows.Close()

	data := newVocabularyData(start, end, nSamples)
	for rows.Next() {
		var interval, total, i int
		if err := rows.Scan(&interval, &total, &i); err != nil {
			return VocabularyData{}, err
		}
		data.add(interval, i, total)
	}
	return data, nil
}
