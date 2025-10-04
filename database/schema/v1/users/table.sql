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