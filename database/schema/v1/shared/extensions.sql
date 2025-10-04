-- Active: 1758835250687@@127.0.0.1@5432@local@fatch
-- For gen_random_uuid()
-- Password hashing.
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- For case-insensitive uniqueness.
CREATE EXTENSION IF NOT EXISTS citext;