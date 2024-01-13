CREATE TABLE apps (
	name	TEXT NOT NULL PRIMARY KEY
);

CREATE TABLE networks (
	name				TEXT NOT NULL PRIMARY KEY,
	internal_network_id	TEXT NOT NULL
);
