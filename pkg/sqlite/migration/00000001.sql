CREATE TABLE configurations (
  k TEXT PRIMARY KEY,
  v TEXT NOT NULL
);

INSERT INTO configurations (k, v) VALUES ('name', 'DeployKit'), ('setup_complete', 'false');
