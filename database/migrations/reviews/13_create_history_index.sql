-- Copyright (c) 2022 Levi Gruspe
-- License: MIT, or AGPLv3 or later

-- +goose Up
CREATE INDEX index_history_reviewed ON history (reviewed);

-- +goose Down
DROP INDEX index_history_reviewed;
