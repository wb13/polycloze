create table seen (
	sentence int primary key,
	last not null default (current_timestamp),
	counter not null default (0)
);
