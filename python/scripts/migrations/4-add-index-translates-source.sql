begin transaction;
	pragma user_version = 4;

	create index if not exists index_translates_source on translates (source);

	commit;
