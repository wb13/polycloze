---- used by review_scheduler

-- Table of review sessions.
-- Unseen items don't appear here.
CREATE TABLE review (
	item TEXT PRIMARY KEY,
	learned NOT NULL DEFAULT CURRENT_TIMESTAMP,
	reviewed NOT NULL DEFAULT CURRENT_TIMESTAMP,
	interval INTEGER NOT NULL,	-- Nanoseconds to add to time now to get the next due date
	due NOT NULL,								-- Date of next review.
	correct BOOLEAN GENERATED ALWAYS AS (interval > 0) VIRTUAL
);
