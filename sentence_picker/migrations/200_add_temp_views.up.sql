create temp view word_difficulty as
select frequency_class.id as word,
			 frequency_class/(1.0 + coalesce((select level from review_schema.most_recent_review where item = frequency_class.word), 0.0)) as difficulty
from language_schema.frequency_class;

create temp view sentence_difficulty as
select sentence,
			 cast(sum(word_difficulty.difficulty) + count(word) as float)
					/ coalesce(1.0, (select counter from seen where sentence = language_schema.contains.sentence))
					as difficulty
from language_schema.contains join word_difficulty using (word)
group by sentence;
