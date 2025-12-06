-- Active: 1758835250687@@127.0.0.1@5432@local@fatch
CREATE TABLE currencies (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE, -- e.g. "US Dollar"
    symbol VARCHAR(5) NOT NULL UNIQUE -- e.g. "USD", "GHS"
);