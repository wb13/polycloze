-- Table of intervals.
create table interval (
	interval primary key,
	correct integer not null default 0,
	incorrect integer not null default 0
);

insert into interval (interval) values (0), (86400000000000);
