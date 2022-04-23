BEGIN TRANSACTION;
PRAGMA user_version = 1;
PRAGMA foreign_keys = ON;

CREATE TABLE Word (
	word PRIMARY KEY,
	frequency NOT NULL
);

CREATE TABLE Sentence (
	id INTEGER PRIMARY KEY,
	text UNIQUE NOT NULL,
	tokenization NOT NULL	-- json array of tokens
);

CREATE TABLE Contains (
	sentence NOT NULL REFERENCES Sentence,
	word NOT NULL REFERENCES Word
);

COMMIT;
