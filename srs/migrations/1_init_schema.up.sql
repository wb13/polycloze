PRAGMA user_version = 1;

-- Append-only table of coefficients, which are used to update intervals after
-- reviewing.
CREATE TABLE Coefficient (
	id INTEGER PRIMARY KEY,
	level INTEGER,
	coefficient FLOAT NOT NULL DEFAULT 2.0
	-- Multiplier used at current level to get the next interval.
	-- See review.go/nextReview for the exceptions.
);

-- Table of review sessions.
-- Unseen items don't appear here.
CREATE TABLE Review (
	id INTEGER PRIMARY KEY,
	word TEXT NOT NULL,
	reviewed NOT NULL DEFAULT CURRENT_TIMESTAMP,
	correct BOOLEAN NOT NULL,		-- Result of this review.

	interval INTEGER NOT NULL,	-- Nanoseconds to add to time now to get the next due date
	due NOT NULL,								-- Date of next review.
	level INTEGER GENERATED ALWAYS AS (
		CAST(floor(log2(2 * interval / 86400000000000 + 1)) AS INTEGER)
	) STORED
);

CREATE TRIGGER insert_default_coefficient BEFORE INSERT ON Review
WHEN NEW.level NOT IN (SELECT level FROM Coefficient)
	BEGIN
		INSERT INTO Coefficient (level, coefficient) VALUES (NEW.level, 2.0);
	END;

-- See https://sqlite.org/lang_select.html#bare_columns_in_an_aggregate_query.
CREATE VIEW MostRecentReview AS
SELECT Review.* FROM Review JOIN (SELECT max(id) AS id FROM Review GROUP BY word)
USING (id);

CREATE VIEW UpdatedCoefficient AS
SELECT Coefficient.* FROM Coefficient JOIN (SELECT max(id) AS id FROM Coefficient GROUP BY level)
USING (id);
