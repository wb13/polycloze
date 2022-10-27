-- Copyright (c) 2022 Levi Gruspe
-- License: MIT, or AGPLv3 or later

-- +goose Up
CREATE TABLE student (
	key PRIMARY KEY CHECK (key = 'me'),
	frequency_class INTEGER NOT NULL DEFAULT 0,
	correct INTEGER NOT NULL DEFAULT 0,	-- only for words seen for the first time
	incorrect INTEGER NOT NULL DEFAULT 0
);

INSERT OR IGNORE INTO student (key) VALUES ('me');

-- +goose Down
DROP TABLE student;
