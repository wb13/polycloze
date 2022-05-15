.read database/migrations/1_init_schema.up.sql

create temp view word_difficulty as
select frequency_class.id as word,
			 frequency_class/(1.0 + coalesce(level, 0.0)) as difficulty
from l2.frequency_class left join most_recent_review on (frequency_class.word = most_recent_review.item);

attach database './cmd/spa.db' as l2;
attach database './cmd/eng.db' as l1;
attach database './cmd/translations.db' as translation;

.timer on
