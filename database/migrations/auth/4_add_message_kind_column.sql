-- Copyright (c) 2022 Levi Gruspe
-- License: MIT, or AGPLv3 or later

-- +goose Up

ALTER TABLE message
ADD COLUMN kind TEXT CHECK (kind IN ('info', 'warning', 'error', 'success'));

-- +goose Down
ALTER TABLE message DROP COLUMN kind;
