-- Copyright (c) 2022 Levi Gruspe
-- License: MIT, or AGPLv3 or later

---- Used by review_scheduler.

-- +goose Up

-- Table of review sessions.
-- Unseen items don't appear here.
CREATE TABLE review (
	item TEXT PRIMARY KEY,
	learned NOT NULL DEFAULT CURRENT_TIMESTAMP,
	reviewed NOT NULL DEFAULT CURRENT_TIMESTAMP,
	interval INTEGER NOT NULL,	-- # of seconds to add to current timestamp to get the next due date
	due NOT NULL,								-- Date of next review
	correct BOOLEAN GENERATED ALWAYS AS (interval > 0) VIRTUAL
);

-- +goose Down
DROP TABLE review;
