BEGIN TRANSACTION;
PRAGMA user_version = 1;

-- Table of coefficients, used to update intervals after reviewing.
CREATE TABLE Coefficient (
	streak INTEGER PRIMARY KEY,
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

CREATE VIEW MostRecentReview AS
SELECT DISTINCT A.word, A.due, A.interval, A.reviewed, A.correct
FROM Review AS A CROSS JOIN Review AS B
WHERE A.word = B.word AND A.reviewed >= B.reviewed;

COMMIT;
