-- Copyright (c) 2022 Levi Gruspe
-- License: MIT, or AGPLv3 or later

-- +goose Up
-- +goose StatementBegin

-- This should only contain one entry at most.
CREATE TABLE vocabulary_size (
	id TEXT PRIMARY KEY DEFAULT 'vocabulary-size' CHECK (id = 'vocabulary-size'),
	t INTEGER NOT NULL DEFAULT (unixepoch('now')),
	v INTEGER NOT NULL
);

CREATE TABLE vocabulary_size_history (
	id INTEGER PRIMARY KEY,
	t INTEGER NOT NULL DEFAULT (unixepoch('now')),
	v INTEGER NOT NULL
);

CREATE TRIGGER trigger_vocabulary_size_after_insert_on_history_case_decrease
AFTER INSERT ON history
FOR EACH ROW
	WHEN coalesce(NEW.interval_before, 0) > 0 AND NEW.interval_after < NEW.interval_before
		BEGIN
			INSERT INTO vocabulary_size (t, v)
			VALUES (NEW.reviewed, 0)
			ON CONFLICT DO UPDATE SET
				t = excluded.t,
				v = max(v - 1, 0);
		END;


CREATE TRIGGER trigger_vocabulary_size_after_insert_on_history_case_increase
AFTER INSERT ON history
FOR EACH ROW
	WHEN coalesce(NEW.interval_before, 0) <= 0 AND coalesce(NEW.interval_before, 0) < NEW.interval_after
		BEGIN
			INSERT INTO vocabulary_size (t, v)
			VALUES (NEW.reviewed, 1)
			ON CONFLICT DO UPDATE SET
				t = excluded.t,
				v = v + 1;
		END;

CREATE TRIGGER trigger_vocabulary_size_history_after_insert_on_vocabulary_size
AFTER INSERT ON vocabulary_size
FOR EACH ROW
	BEGIN
		INSERT INTO vocabulary_size_history (t, v) VALUES (NEW.t, NEW.v);
	END;

CREATE TRIGGER trigger_vocabulary_size_history_after_update_on_vocabulary_size
AFTER UPDATE ON vocabulary_size
FOR EACH ROW
	BEGIN
		INSERT INTO vocabulary_size_history (t, v) VALUES (NEW.t, NEW.v);
	END;

-- +goose StatementEnd

-- +goose Down

DROP TRIGGER trigger_vocabulary_size_history_after_insert_on_vocabulary_size;
DROP TRIGGER trigger_vocabulary_size_history_after_update_on_vocabulary_size;

DROP TRIGGER trigger_vocabulary_size_after_insert_on_history_case_decrease;
DROP TRIGGER trigger_vocabulary_size_after_insert_on_history_case_increase;

DROP TABLE vocabulary_size_history;
DROP TABLE vocabulary_size;
