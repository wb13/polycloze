-- Copyright (c) 2022 Levi Gruspe
-- License: MIT, or AGPLv3 or later

-- +goose Up
CREATE TABLE csrf_token (
	session_id TEXT NOT NULL REFERENCES user_session,
	token TEXT NOT NULL
);

-- +goose Down
DROP TABLE csrf_token;
