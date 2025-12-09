DROP TABLE IF EXISTS currencies CASCADE;
CREATE TABLE currencies (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE, -- e.g. "US Dollar"
    symbol VARCHAR(5) NOT NULL UNIQUE -- e.g. "USD", "GHS"
);
