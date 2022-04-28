BEGIN TRANSACTION;
PRAGMA user_version = 0;

DROP VIEW MostRecentReview;
DROP TABLE Review;

COMMIT;
