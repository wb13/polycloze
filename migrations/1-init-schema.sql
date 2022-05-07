begin transaction;

	pragma user_version = 1;

	-- language metadata (only has one column)
	create table info (
		id integer primary key default 0,
		code not null,
		language not null,	-- only one language per database
		check (id = 0)
	);

	create table word (
		id integer primary key,
		word unique not null,
		frequency not null
	);

	create table sentence (
		id integer primary key,
		tatoeba_id integer unique,	-- null for non-tatoeba sentences
		text unique not null,
		tokens not null	-- json array of strings
	);

	create table contains (
		sentence not null references sentence,
		word not null references word
	);

	commit;

	-- TODO
	-- select id, word, cast(floor(0.5 - log2(frequency/98578.0)) as int) as frequency_class from word;
