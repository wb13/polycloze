-- Copyright (c) 2022 Levi Gruspe
-- License: MIT, or AGPLv3 or later

-- +goose Up
-- +goose StatementBegin

-- Keeps track of # of words at each interval.
CREATE TABLE interval_count (
	interval INTEGER, 			-- Interval size
	total INTEGER NOT NULL,	-- Cumulative sum
	t INTEGER NOT NULL DEFAULT (unixepoch('now'))
);

CREATE VIEW most_recent_interval_count AS
SELECT interval, total, max(t) AS t
FROM interval_count
GROUP BY interval;

CREATE TRIGGER interval_count_trigger_after_insert_review
AFTER INSERT ON review
FOR EACH ROW
	BEGIN
		-- Increase new interval.
		INSERT INTO interval_count (interval, total)
		VALUES (
			NEW.interval,
			coalesce(
				(
					SELECT total + 1
					FROM most_recent_interval_count
					WHERE interval = NEW.interval
				),
				1
			)
		);
	END;

CREATE TRIGGER interval_count_trigger_after_update_review
AFTER UPDATE OF interval ON review
FOR EACH ROW
	BEGIN
		-- Increase new interval, decrease old interval.
		INSERT INTO interval_count (interval, total) VALUES
		(
			NEW.interval,
			coalesce(
				(
					SELECT total + 1
					FROM most_recent_interval_count
					WHERE interval = NEW.interval
				),
				1
			)
		),
		(
			OLD.interval,
			coalesce(
				(
					SELECT total - 1
					FROM most_recent_interval_count
					WHERE interval = OLD.interval
				),
				0
			)
		);
	END;

CREATE TRIGGER interval_count_trigger_after_delete_review
AFTER DELETE ON review
FOR EACH ROW
	BEGIN
		-- Decrease old interval.
		INSERT INTO interval_count (interval, total)
		VALUES (
			OLD.interval,
			coalesce(
				(
					SELECT total - 1
					FROM most_recent_interval_count
					WHERE interval = OLD.interval
				),
				0
			)
		);
	END;

-- +goose StatementEnd

-- +goose Down

DROP TRIGGER interval_count_trigger_after_insert_review;
DROP TRIGGER interval_count_trigger_after_update_review;
DROP TRIGGER interval_count_trigger_after_delete_review;
DROP VIEW most_recent_interval_count;
DROP TABLE interval_count;
