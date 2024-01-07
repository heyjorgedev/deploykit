CREATE TABLE apps (
	id		INTEGER PRIMARY KEY AUTOINCREMENT,
	name	TEXT NOT NULL UNIQUE,
	network	TEXT NOT NULL
);

CREATE TABLE networks (
	name	TEXT NOT NULL PRIMARY KEY,
	network_id TEXT NOT NULL UNIQUE
);

