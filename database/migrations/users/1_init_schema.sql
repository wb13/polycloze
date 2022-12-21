-- Copyright (c) 2022 Levi Gruspe
-- License: MIT, or AGPLv3 or later

-- +goose Up
-- Key-value store of user data.
CREATE TABLE user_data (
	name TEXT PRIMARY KEY,
	value TEXT
);

-- +goose Down
DROP TABLE user_data;
