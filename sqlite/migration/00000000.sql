CREATE TABLE sessions (
  token TEXT PRIMARY KEY,
  data BLOB NOT NULL,
  expiry REAL NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions(expiry);

CREATE TABLE users (
  id INTEGER PRIMARY KEY,
  username TEXT NOT NULL,
  password TEXT NOT NULL,
  email TEXT NOT NULL
);
