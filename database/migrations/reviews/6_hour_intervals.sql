-- Copyright (c) 2022 Levi Gruspe
-- License: MIT, or AGPLv3 or later

-- +goose Up

-- Reduce interval resolution from seconds to hours.
-- Uses 12-step migration: https://sqlite.org/lang_altertable.html#making_other_kinds_of_table_schema_changes

CREATE TABLE new_interval (
	interval PRIMARY KEY,	-- In hours
	correct INTEGER NOT NULL DEFAULT 0,
	Incorrect INTEGER NOT NULL DEFAULT 0
);

INSERT INTO new_interval
SELECT CAST(round(interval / 3600.0) AS INTEGER) AS x, sum(correct), sum(incorrect)
FROM interval
GROUP BY x
ORDER BY x ASC;

DROP TABLE interval;

ALTER TABLE new_interval RENAME TO interval;

INSERT OR IGNORE INTO interval (interval) VALUES (0);
UPDATE review SET interval = CAST(round(interval / 3600.0) AS INTEGER);

-- +goose Down

UPDATE interval SET interval = interval * 3600;
UPDATE review SET interval = interval * 3600;
