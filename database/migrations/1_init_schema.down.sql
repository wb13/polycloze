-- used by review_scheduler
drop view most_recent_review;
drop view updated_coefficient;

drop trigger insert_default_coefficient;

drop table review;
drop table coefficient;

-- used by sentence_picker
drop table seen;
