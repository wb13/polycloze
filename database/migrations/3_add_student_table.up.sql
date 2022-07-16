create table student (
	key primary key check (key = 'me'),
	frequency_class integer not null default 0,
	correct integer not null default 0,	-- only for words seen for the first time
	incorrect integer not null default 0
);

insert or ignore into student (key) values ('me');
