CREATE TABLE signatures (
    signature TEXT NOT NULL,
    hash1 TEXT PRIMARY KEY,
    hash2 INTEGER NOT NULL,
    hash3 INTEGER NOT NULL
);

CREATE TABLE documents (
    id TEXT PRIMARY KEY,
    hash1 TEXT NOT NULL
);
