-- Copyright (c) 2022 Levi Gruspe
-- License: MIT, or AGPLv3 or later

-- +goose Up
-- +goose StatementBegin

CREATE TABLE estimated_level (
	id TEXT PRIMARY KEY DEFAULT 'estimated-level' CHECK (id = 'estimated-level'),
	t INTEGER NOT NULL DEFAULT (unixepoch('now')),
	v INTEGER NOT NULL DEFAULT 0,
	correct INTEGER NOT NULL DEFAULT 0,
	incorrect INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE estimated_level_history (
	id INTEGER PRIMARY KEY,
	t INTEGER NOT NULL DEFAULT (unixepoch('now')),
	v INTEGER NOT NULL
);

CREATE TRIGGER trigger_estimated_level_history_after_insert_on_estimated_level
AFTER INSERT ON estimated_level
FOR EACH ROW
	BEGIN
		INSERT INTO estimated_level_history (t, v) VALUES (NEW.t, NEW.v);
	END;

CREATE TRIGGER trigger_estimated_level_history_after_update_on_estimated_level
AFTER UPDATE ON estimated_level
FOR EACH ROW
	WHEN OLD.v != NEW.v
		BEGIN
			INSERT INTO estimated_level_history (t, v) VALUES (NEW.t, NEW.v);
		END;

-- +goose StatementEnd

-- +goose Down

DROP TRIGGER trigger_estimated_level_history_after_insert_on_estimated_level;
DROP TRIGGER trigger_estimated_level_history_after_update_on_estimated_level;

DROP TABLE estimated_level_history;
DROP TABLE estimated_level;
