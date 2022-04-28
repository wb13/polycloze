BEGIN TRANSACTION;
PRAGMA user_version = 1;

CREATE TABLE Review (
	word TEXT NOT NULL,
	reviewed NOT NULL DEFAULT CURRENT_TIMESTAMP,

	interval INTEGER NOT NULL,
	due NOT NULL DEFAULT,
	correct BOOLEAN NOT NULL
);

CREATE VIEW MostRecentReview AS
SELECT DISTINCT A.word, A.due, A.interval, A.reviewed, A.correct
FROM Review AS A CROSS JOIN Review AS B
WHERE A.word = B.word AND A.reviewed >= B.reviewed;

COMMIT;
