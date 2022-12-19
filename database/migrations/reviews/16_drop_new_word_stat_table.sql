-- Copyright (c) 2022 Levi Gruspe
-- License: MIT, or AGPLv3 or later

-- +goose Up
DROP TABLE new_word_stat;

-- +goose Down
-- See `4_replace_student_table.sql`.
CREATE TABLE new_word_stat (
	frequency_class INTEGER PRIMARY KEY,
	correct INTEGER NOT NULL DEFAULT 0,
	incorrect INTEGER NOT NULL DEFAULT 0
);
