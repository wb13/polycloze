-- Copyright (c) 2022 Levi Gruspe
-- License: MIT, or AGPLv3 or later

-- +goose Up

DROP VIEW most_recent_interval_count;

CREATE VIEW most_recent_interval_count AS
SELECT interval, total, t FROM (
	SELECT interval, total, t, max(ROWID)
	FROM interval_count
	GROUP BY interval
);


-- +goose Down

DROP VIEW most_recent_interval_count;

CREATE VIEW most_recent_interval_count AS
SELECT interval, total, max(t) AS t
FROM interval_count
GROUP BY interval;
