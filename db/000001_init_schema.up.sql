CREATE TABLE IF NOT EXISTS persons(
    id TEXT PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    age INT NOT NULL,
    works BOOLEAN NOT NULL,
    password TEXT NOT NULL,
    refreshtoken TEXT
    );