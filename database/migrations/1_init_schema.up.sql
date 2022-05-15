---- used by review_scheduler

-- Append-only table of coefficients, which are used to update intervals after
-- reviewing.
CREATE TABLE coefficient (
	id INTEGER PRIMARY KEY,
	level INTEGER,
	coefficient FLOAT NOT NULL DEFAULT 2.0
	-- Multiplier used at current level to get the next interval.
	-- See review.go/nextReview for the exceptions.
);

-- Table of review sessions.
-- Unseen items don't appear here.
CREATE TABLE review (
	id INTEGER PRIMARY KEY,
	item TEXT NOT NULL,
	reviewed NOT NULL DEFAULT CURRENT_TIMESTAMP,
	interval INTEGER NOT NULL,	-- Nanoseconds to add to time now to get the next due date
	due NOT NULL,								-- Date of next review.

	correct BOOLEAN GENERATED ALWAYS AS (interval > 0) VIRTUAL,
	level INTEGER GENERATED ALWAYS AS (
		CAST(floor(log2(2 * interval / 86400000000000 + 1)) AS INTEGER)
	) STORED
);

CREATE TRIGGER insert_default_coefficient BEFORE INSERT ON Review
WHEN NEW.level NOT IN (SELECT level FROM coefficient)
	BEGIN
		INSERT INTO coefficient (level, coefficient) VALUES (NEW.level, 2.0);
	END;

-- See https://sqlite.org/lang_select.html#bare_columns_in_an_aggregate_query.
CREATE VIEW most_recent_review AS
SELECT review.* FROM review JOIN (SELECT max(id) AS id FROM review GROUP BY item)
USING (id);

CREATE VIEW updated_coefficient AS
SELECT coefficient.* FROM coefficient JOIN (SELECT max(id) AS id FROM coefficient GROUP BY level)
USING (id);
