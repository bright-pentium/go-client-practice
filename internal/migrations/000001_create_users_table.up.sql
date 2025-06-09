-- users table
CREATE TABLE users (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    account TEXT NOT NULL UNIQUE,
    password_hash BYTEA NOT NULL
);

-- clients table
CREATE TABLE clients (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    scope TEXT[] NOT NULL,
    secret_hash BYTEA NOT NULL
);
