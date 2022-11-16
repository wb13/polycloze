-- Copyright (c) 2022 Levi Gruspe
-- License: MIT, or AGPLv3 or later

-- +goose Up

-- Not reversible.

DROP TRIGGER interval_count_trigger_after_insert_review;
DROP TRIGGER interval_count_trigger_after_update_review;
DROP TRIGGER interval_count_trigger_after_delete_review;
DROP VIEW most_recent_interval_count;
DROP TABLE interval_count;
