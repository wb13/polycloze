-- Table of intervals.
create table interval (
	interval primary key,
	correct integer not null default 0,
	incorrect integer not null default 0
);
