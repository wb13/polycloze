begin transaction;
	pragma user_version = 5;

	create index if not exists index_word_frequency_class on word (frequency_class);

	commit;
