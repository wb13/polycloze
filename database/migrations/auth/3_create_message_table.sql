-- Copyright (c) 2022 Levi Gruspe
-- License: MIT, or AGPLv3 or later

-- +goose Up
CREATE TABLE message (
	id INTEGER PRIMARY KEY,
	session_id TEXT REFERENCES user_session,
	created INTEGER NOT NULL DEFAULT (unixepoch('now')),
	message TEXT NOT NULL
);

CREATE INDEX index_message_session_id ON message (session_id);

-- +goose Down
DROP INDEX index_message_session_id;
DROP TABLE message;
