-- Copyright (c) 2022 Levi Gruspe
-- License: MIT, or AGPLv3 or later

-- +goose Up

-- Generate `review.due` as a virtual column instead of storing it.

ALTER TABLE review DROP COLUMN due;
ALTER TABLE review ADD COLUMN due INTEGER NOT NULL GENERATED ALWAYS AS (reviewed + 3600*interval) VIRTUAL;


-- +goose Down

ALTER TABLE review DROP COLUMN due;
ALTER TABLE review ADD COLUMN due INTEGER NOT NULL DEFAULT (reviewed + 3600*interval);
