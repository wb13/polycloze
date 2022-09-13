-- Generated using python/scripts/dump.py
BEGIN TRANSACTION;
CREATE TABLE contains (
		sentence integer not null references sentence,
		word integer not null references word
		);
CREATE TABLE language (
		id text primary key check (id = 'l1' or id = 'l2'),
		code char(3) not null check (length(code) = 3),
		name text not null
		, bcp47 text not null);
CREATE TABLE sentence (
		id integer primary key,
		tatoeba_id integer unique,	-- null for non-tatoeba sentences
		text text unique not null,
		tokens text not null,	-- json array of strings
		frequency_class integer not null	-- max frequency_class among all words in sentence
		);
CREATE TABLE translates (
		source integer not null,	-- references sentence.tatoeba_id
		target integer not null		-- references translation.tatoeba_id
		);
CREATE TABLE translation (
		id integer primary key,
		tatoeba_id integer unique,	-- null for non-tatoeba sentences
		text text unique not null
		);
CREATE TABLE word (
		id integer primary key,
		word text unique not null,
		frequency_class integer not null
		);
CREATE INDEX index_contains_word on contains (word);
COMMIT;
