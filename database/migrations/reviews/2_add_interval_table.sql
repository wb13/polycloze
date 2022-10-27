-- Copyright (c) 2022 Levi Gruspe
-- License: MIT, or AGPLv3 or later

-- +goose Up
CREATE TABLE interval (
	interval PRIMARY KEY,	-- In seconds
	correct INTEGER NOT NULL DEFAULT 0,
	incorrect INTEGER NOT NULL DEFAULT 0
);

INSERT INTO interval (interval) VALUES (0);

-- +goose Down
DROP TABLE interval;
