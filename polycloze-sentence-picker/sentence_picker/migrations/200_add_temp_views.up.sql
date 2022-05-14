create temp view word_difficulty as
select frequency_class.id as word,
			 frequency_class/(1.0 + coalesce(level, 0.0)) as difficulty
from frequency_class left join most_recent_review on (frequency_class.word = most_recent_review.item);

create temp view sentence_difficulty as
select sentence,
			 sum(word_difficulty.difficulty) + count(word)
					+ coalesce((select counter from seen where sentence = contains.sentence), 0.0)
					as difficulty
from contains join word_difficulty using (word)
group by sentence;
