PRAGMA user_version = 1;

-- Append-only table of coefficients, used to update intervals after reviewing.
CREATE TABLE Coefficient (
	streak INTEGER,

	-- Multiplier to get the interval for the next level, except for going from
	-- level 0 to 1 (sets due date to the next day instead)
	coefficient FLOAT NOT NULL DEFAULT 2.0,
	updated NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Table of review sessions.
-- Unseen items (newly learned words) won't appear here.
CREATE TABLE Review (
	word TEXT NOT NULL,
	reviewed NOT NULL DEFAULT CURRENT_TIMESTAMP,
	correct BOOLEAN NOT NULL,		-- Result of this review.

	interval INTEGER NOT NULL,	-- Used to compute next due date.
	due NOT NULL,								-- Date of next review.

	streak INTEGER REFERENCES Coefficient
);

CREATE TRIGGER insert_default_coefficient BEFORE INSERT ON Review
WHEN NEW.streak NOT IN (SELECT streak FROM Coefficient)
	BEGIN
		INSERT INTO Coefficient (streak, coefficient) VALUES (NEW.streak, 2.0);
	END;

-- See https://sqlite.org/lang_select.html#bare_columns_in_an_aggregate_query.
CREATE VIEW MostRecentReview AS
SELECT Review.* FROM Review JOIN (SELECT max(rowid) FROM Review GROUP BY word)
USING (rowid);

CREATE VIEW UpdatedCoefficient AS
SELECT Coefficient.* FROM Coefficient JOIN (SELECT max(rowid) FROM Coefficient GROUP BY streak)
USING (rowid);
