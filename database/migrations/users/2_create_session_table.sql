-- Copyright (c) 2022 Levi Gruspe
-- License: MIT, or AGPLv3 or later

-- +goose Up
CREATE TABLE user_session (
	session_id TEXT PRIMARY KEY CHECK(session_id != ''),
	user_id INTEGER REFERENCES user,	-- null if user is not logged in
	username TEXT 										-- null if user is not logged in
);

-- +goose Down
DROP TABLE user_session;
