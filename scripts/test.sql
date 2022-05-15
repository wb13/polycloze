.read database/migrations/1_init_schema.up.sql

create temp view word_difficulty as
select word.id as word,
			 frequency_class/(1.0 + coalesce(level, 0.0)) as difficulty
from l2.word left join most_recent_review on (word = most_recent_review.item);

attach database './cmd/spa.db' as l2;
attach database './cmd/eng.db' as l1;
attach database './cmd/translations.db' as translation;

.timer on
