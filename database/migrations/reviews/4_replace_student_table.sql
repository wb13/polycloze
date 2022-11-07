-- Copyright (c) 2022 Levi Gruspe
-- License: MIT, or AGPLv3 or later

-- +goose Up
DROP TABLE student;

-- Statistics on newly seen words.
CREATE TABLE new_word_stat (
	frequency_class INTEGER PRIMARY KEY,
	correct INTEGER NOT NULL DEFAULT 0,
	incorrect INTEGER NOT NULL DEFAULT 0
);

-- +goose Down
-- See `3_add_student_table.sql` upgrade.
CREATE TABLE student (
	key PRIMARY KEY CHECK (key = 'me'),
	frequency_class INTEGER NOT NULL DEFAULT 0,
	correct INTEGER NOT NULL DEFAULT 0,	-- only for words seen for the first time
	incorrect INTEGER NOT NULL DEFAULT 0
);

INSERT OR IGNORE INTO student (key) VALUES ('me');


