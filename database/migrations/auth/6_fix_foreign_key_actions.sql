-- Copyright (c) 2022 Levi Gruspe
-- License: MIT, or AGPLv3 or later

-- +goose NO TRANSACTION
-- +goose Up
-- +goose StatementBegin

-- Fixes foreign keys in `2_create_session_table.sql` and
-- `3_create_message_table.sql`.
-- Follows the 12-step schema change in:
-- https://www.sqlite.org/lang_altertable.html#making_other_kinds_of_table_schema_changes

-- Fix `user_session` schema.
PRAGMA foreign_keys = OFF;

BEGIN TRANSACTION;

CREATE TABLE new_user_session (
	session_id TEXT PRIMARY KEY CHECK(session_id != ''),
	created INTEGER NOT NULL DEFAULT (unixepoch('now')),
	updated INTEGER NOT NULL DEFAULT (unixepoch('now')),

	-- The following columns are null if the user is not logged in.
	user_id INTEGER REFERENCES user ON DELETE CASCADE,
	username TEXT
);

INSERT INTO new_user_session (session_id, created, updated, user_id, username)
SELECT session_id, created, updated, user_id, username
FROM user_session;

DROP TABLE user_session;

ALTER TABLE new_user_session RENAME TO user_session;

PRAGMA foreign_key_check;

-- Fix `message` schema.

DROP INDEX index_message_session_id;

CREATE TABLE new_message (
	id INTEGER PRIMARY KEY,
	session_id TEXT REFERENCES user_session ON DELETE CASCADE,
	created INTEGER NOT NULL DEFAULT (unixepoch('now')),
	message TEXT NOT NULL,
	kind TEXT CHECK (kind IN ('info', 'warning', 'error', 'success')),
	context TEXT DEFAULT ''
);

INSERT INTO new_message (id, session_id, created, message, kind, context)
SELECT id, session_id, created, message, kind, context
FROM message;

DROP TABLE message;

ALTER TABLE new_message RENAME TO message;

CREATE INDEX index_message_session_id ON message (session_id);

PRAGMA foreign_key_check;

END;

PRAGMA foreign_keys = ON;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

PRAGMA foreign_keys = OFF;

BEGIN TRANSACTION;

-- Undo changes to `message` schema.

DROP INDEX index_message_session_id;

-- Summary of changes from version 3 to 5.
CREATE TABLE old_message (
	id INTEGER PRIMARY KEY,
	session_id TEXT REFERENCES user_session,
	created INTEGER NOT NULL DEFAULT (unixepoch('now')),
	message TEXT NOT NULL,
	kind TEXT CHECK (kind IN ('info', 'warning', 'error', 'success')),
	context TEXT DEFAULT ''
);

INSERT INTO old_message (id, session_id, created, message, kind, context)
SELECT id, session_id, created, message, kind, context
FROM message;

DROP TABLE message;

ALTER TABLE old_message RENAME TO message;

CREATE INDEX index_message_session_id ON message (session_id);

PRAGMA foreign_key_check;

-- Undo changes to `user_session` schema.

CREATE TABLE old_user_session (
	session_id TEXT PRIMARY KEY CHECK(session_id != ''),
	created INTEGER NOT NULL DEFAULT (unixepoch('now')),
	updated INTEGER NOT NULL DEFAULT (unixepoch('now')),
	user_id INTEGER REFERENCES user,	-- null if user is not logged in
	username TEXT 										-- null if user is not logged in
);

INSERT INTO old_user_session (session_id, created, updated, user_id, username)
SELECT session_id, created, updated, user_id, username
FROM user_session;

DROP TABLE user_session;

ALTER TABLE old_user_session RENAME TO user_session;

PRAGMA foreign_key_check;

END;

PRAGMA foreign_keys = ON;

-- +goose StatementEnd
