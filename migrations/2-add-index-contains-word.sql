begin transaction;
	pragma user_version = 2;

	create index if not exists index_contains_word on contains (word);

	commit;
