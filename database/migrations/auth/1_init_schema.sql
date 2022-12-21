-- Copyright (c) 2022 Levi Gruspe
-- License: MIT, or AGPLv3 or later

-- +goose Up
CREATE TABLE user (
	id INTEGER PRIMARY KEY,
	username TEXT UNIQUE NOT NULL CHECK(username != ''),

	-- Salted and hashed using bcrypt.
	password TEXT NOT NULL CHECK(password != '')
);

-- +goose Down
DROP TABLE user;
