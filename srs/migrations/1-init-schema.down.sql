BEGIN TRANSACTION;
PRAGMA user_version = 0;

DROP VIEW MostRecentReview;

DROP TRIGGER insert_default_coefficient;

DROP TABLE Review;
DROP TABLE Coefficient;

COMMIT;
