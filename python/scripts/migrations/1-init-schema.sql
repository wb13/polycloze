begin transaction;
	pragma user_version = 1;

	create table if not exists language (
		id text primary key check (id = 'l1' or id = 'l2'),
		code char(3) not null check (length(code) = 3),
		name text not null
		);

	create table if not exists word (
		id integer primary key,
		word text unique not null,
		frequency_class integer not null
		);

	create table if not exists sentence (
		id integer primary key,
		tatoeba_id integer unique,	-- null for non-tatoeba sentences
		text text unique not null,
		tokens text not null,	-- json array of strings
		frequency_class integer not null	-- max frequency_class among all words in sentence
		);

	create table if not exists contains (
		sentence integer not null references sentence,
		word integer not null references word
		);

	create table if not exists translation (
		id integer primary key,
		tatoeba_id integer unique,	-- null for non-tatoeba sentences
		text text unique not null
		);

	create table if not exists translates (
		source integer not null,	-- references sentence.tatoeba_id
		target integer not null		-- references translation.tatoeba_id
		);

	commit;
