begin transaction;
	pragma user_version = 1;

	create table word (
		id integer primary key,
		word unique not null,
		frequency integer not null,
		frequency_class integer
	);

	create table sentence (
		id integer primary key,
		tatoeba_id integer unique,	-- null for non-tatoeba sentences
		text unique not null,
		tokens not null,	-- json array of strings
		frequency_class integer	-- max frequency_class among all words in sentence
	);

	create table contains (
		sentence not null references sentence,
		word not null references word
	);

	create index index_contains_word on contains (word);

	commit;
