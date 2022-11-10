-- Copyright (c) 2022 Levi Gruspe
-- License: MIT, or AGPLv3 or later

-- +goose Up

-- Alter review table to use UNIX timestamps.
-- 12 step schema change (see https://sqlite.org/lang_altertable.html#making_other_kinds_of_table_schema_changes):

-- New schema of review
CREATE TABLE new_review (
		item TEXT PRIMARY KEY,
		learned INTEGER NOT NULL DEFAULT (unixepoch('now')),
		reviewed INTEGER NOT NULL DEFAULT (unixepoch('now')),
		interval INTEGER NOT NULL,	-- # of seconds from reviewed to due
		due INTEGER NOT NULL,				-- Date of next review (UNIX timestamp)
		correct BOOLEAN GENERATED ALWAYS AS (interval > 0) VIRTUAL
);

-- NOTE `review.interval` does not reference the `interval` table,
-- because merging intervals would require all reviews with interval ID
-- >= the merged interval ID to be updated if `review.interval` references
-- the `interval` table's key.

-- Copy data from old to new table.
INSERT INTO new_review
SELECT item, unixepoch(learned), unixepoch(reviewed), interval, unixepoch(due)
FROM review;

-- Drop the old table.
DROP TABLE review;

-- Rename the new table.
ALTER TABLE new_review RENAME TO review;


-- +goose Down

-- Old review schema (same as in `1_init_schema.sql`)
CREATE TABLE old_review (
	item TEXT PRIMARY KEY,
	learned NOT NULL DEFAULT CURRENT_TIMESTAMP,
	reviewed NOT NULL DEFAULT CURRENT_TIMESTAMP,
	interval INTEGER NOT NULL,	-- # of seconds to add to current timestamp to get the next due date
	due NOT NULL,								-- Date of next review
	correct BOOLEAN GENERATED ALWAYS AS (interval > 0) VIRTUAL
);

INSERT INTO old_review
SELECT item, datetime(learned, 'unixepoch'), datetime(reviewed, 'unixepoch'), interval, datetime(due, 'unixepoch')
FROM review;

DROP TABLE review;

ALTER TABLE old_review RENAME TO review;
