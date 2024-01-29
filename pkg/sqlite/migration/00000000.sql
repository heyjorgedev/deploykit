CREATE TABLE sessions (
  token TEXT PRIMARY KEY,
  data BLOB NOT NULL,
  expiry REAL NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions(expiry);

CREATE TABLE users (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL,
  username TEXT NOT NULL,
  passwordHash TEXT NOT NULL
);

INSERT INTO users (name, username, passwordHash) VALUES ('Admin', 'admin', '$2a$12$M2YjYX9G1JrZvuXBPvjVjOshAXsO.HWMqHCjC1iGuZc4RvkBHt4DW');

