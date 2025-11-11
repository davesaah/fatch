-- Active: 1758835250687@@127.0.0.1@5432@local@fatch
DROP TABLE IF EXISTS users CASCADE;

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    username citext NOT NULL UNIQUE,
    email citext NOT NULL UNIQUE,
    passwd TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- index for username fast lookups
CREATE INDEX idx_users_username ON users (username);

-- index for email fast lookups
CREATE INDEX idx_users_email ON users (email);

-- add updated_at column
ALTER TABLE users
ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT now();
