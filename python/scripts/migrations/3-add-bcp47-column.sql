begin transaction;
	pragma user_version = 3;

	alter table language add column bcp47 text not null;

	commit;
