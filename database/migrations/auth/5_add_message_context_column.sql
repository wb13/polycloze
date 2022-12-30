-- Copyright (c) 2022 Levi Gruspe
-- License: MIT, or AGPLv3 or later

-- +goose Up
ALTER TABLE message ADD COLUMN context TEXT DEFAULT '';

-- +goose Down
ALTER TABLE message DROP COLUMN context;
