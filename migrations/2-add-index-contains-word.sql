begin transaction;
	pragma user_version = 2;

	create index if not exists index_contains_word on contains (word);
	-- create index if not exists index_translates_source on translates (source);

	commit;
