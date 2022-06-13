begin transaction;
	pragma user_version = 1;

	create table translation (
		l1,	-- references <l1.db>.sentence.tatoeba_id
		l2	-- references <l2.db>.sentence.tatoeba_id
	);

	create index index_translation_l2 on translation (l2);

	commit;
