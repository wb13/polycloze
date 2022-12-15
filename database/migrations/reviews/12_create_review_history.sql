-- Copyright (c) 2022 Levi Gruspe
-- License: MIT, or AGPLv3 or later

-- +goose Up
-- +goose StatementBegin

-- Append-only history of reviews.
-- Shouldn't be updated nor deleted from.
CREATE TABLE history (
	word TEXT NOT NULL,
	reviewed INTEGER NOT NULL DEFAULT (unixepoch('now')),
	-- learned = reviewed of first occurrence of item
	interval_before INTEGER,				-- # of hours; before review
	interval_after INTEGER NOT NULL	-- after review
);

CREATE TRIGGER trigger_history_after_insert_on_review
AFTER INSERT ON review
FOR EACH ROW
	BEGIN
		INSERT INTO history (word, reviewed, interval_before, interval_after)
		VALUES (NEW.item, NEW.reviewed, NULL, NEW.interval);
	END;

CREATE TRIGGER trigger_history_after_update_of_reviewed_on_review
AFTER UPDATE OF reviewed ON review
FOR EACH ROW
	BEGIN
		INSERT INTO history (word, reviewed, interval_before, interval_after)
		VALUES (NEW.item, NEW.reviewed, OLD.interval, NEW.interval);
	END;

CREATE TRIGGER trigger_history_after_delete_on_review
AFTER DELETE ON review
FOR EACH ROW
	BEGIN
		-- Delete all entries for the deleted item, so that the review table can be
		-- reconstructed from the review history.
		DELETE FROM history WHERE word = OLD.item;
	END;

-- +goose StatementEnd

-- +goose Down

DROP TRIGGER trigger_history_after_delete_on_review;
DROP TRIGGER trigger_history_after_update_of_reviewed_on_review;
DROP TRIGGER trigger_history_after_insert_on_review;
DROP TABLE history;
