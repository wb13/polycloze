-- Copyright (c) 2022 Levi Gruspe
-- License: MIT, or AGPLv3 or later

-- Table of intervals.
create table interval (
	interval primary key,	-- In seconds
	correct integer not null default 0,
	incorrect integer not null default 0
);

insert into interval (interval) values (0);
