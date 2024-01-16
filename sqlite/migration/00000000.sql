CREATE TABLE projects (
    id      INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	name	TEXT NOT NULL
);

CREATE TABLE networks (
    id                  INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	name				TEXT NOT NULL,
	internal_network_id	TEXT NOT NULL
);
