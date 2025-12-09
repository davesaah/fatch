-- v1 of DB Schema for fatch

CREATE SCHEMA IF NOT EXISTS fatch;
SET search_path TO fatch;

-- For gen_random_uuid()
-- Password hashing.
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- For case-insensitive uniqueness.
CREATE EXTENSION IF NOT EXISTS citext;

----------------------------------------------------------------
-- USERS 
----------------------------------------------------------------
DROP TABLE IF EXISTS users CASCADE;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_email;

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    username citext NOT NULL UNIQUE,
    email citext NOT NULL UNIQUE,
    passwd TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- index for username fast lookups
CREATE INDEX idx_users_username ON users (username);

-- index for email fast lookups
CREATE INDEX idx_users_email ON users (email);

DROP FUNCTION IF EXISTS create_user (citext, citext, TEXT);
CREATE OR REPLACE FUNCTION create_user(
    p_username citext,
    p_email citext,
    p_passwd TEXT
) RETURNS VOID
AS $$
BEGIN
    INSERT INTO users(username, email, passwd)
    VALUES (p_username, p_email, p_passwd);
EXCEPTION
    WHEN unique_violation THEN
        RAISE EXCEPTION 'Username or email already exists';
END;
$$ LANGUAGE plpgsql;

DROP FUNCTION IF EXISTS change_password (UUID, TEXT, TEXT);
CREATE OR REPLACE FUNCTION change_password (
    p_user_id UUID,
    p_old_passwd TEXT,
    p_new_passwd TEXT
) RETURNS VOID
 AS $$
DECLARE
    stored_hash TEXT;
BEGIN
    SELECT passwd INTO stored_hash
    FROM users
    WHERE id = p_user_id;

    IF stored_hash IS NULL THEN
        RAISE EXCEPTION 'User not found';
    END IF;

    IF stored_hash <> crypt(p_old_passwd, stored_hash) THEN
        RAISE EXCEPTION 'Your old password is incorrect';
    END IF;

    UPDATE users
    SET passwd = p_new_passwd
    WHERE id = p_user_id;
END;
$$ LANGUAGE plpgsql;

DROP FUNCTION IF EXISTS verify_password (citext, TEXT);
CREATE OR REPLACE FUNCTION verify_password(
    p_username citext,
    p_passwd TEXT
) RETURNS UUID
AS $$
DECLARE
    stored_hash TEXT;
    user_id UUID;
    is_valid BOOLEAN;
BEGIN
    SELECT passwd, id INTO stored_hash, user_id
    FROM users
    WHERE username = p_username
    LIMIT 1;

    IF stored_hash IS NULL THEN
        RAISE EXCEPTION 'Incorrect credentials';
    ELSE
        is_valid := stored_hash = crypt(p_passwd, stored_hash);
    END IF;

    IF NOT is_valid THEN
        RAISE EXCEPTION 'Incorrect credentials';
    ELSE
        RETURN user_id;
    END IF;
END;
$$ LANGUAGE plpgsql;

DROP FUNCTION IF EXISTS get_user_by_id (uuid);
CREATE OR REPLACE FUNCTION get_user_by_id(p_uuid uuid)
RETURNS TABLE(
    username citext,
    email citext
)
AS $$
BEGIN
    RETURN QUERY SELECT u.username, u.email
    FROM users u
    WHERE id = p_uuid;
END;
$$ LANGUAGE plpgsql;

-- Trigger function to hash passwords using bcrypt
DROP FUNCTION IF EXISTS hash_user_password;
CREATE OR REPLACE FUNCTION hash_user_password()
RETURNS TRIGGER AS $$
BEGIN
    -- Only hash if the password is changed or on insert
    IF NEW.passwd IS DISTINCT FROM OLD.passwd THEN
        NEW.passwd := crypt(NEW.passwd, gen_salt('bf'));
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Attach trigger to users table
DROP TRIGGER IF EXISTS trg_hash_password ON users;
CREATE TRIGGER trg_hash_password
BEFORE INSERT OR UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION hash_user_password();

-- Trigger function to update updated_at
DROP FUNCTION IF EXISTS update_users_timestamp;
CREATE OR REPLACE FUNCTION update_users_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at := now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_users_update_timestamp ON users;
CREATE TRIGGER trg_users_update_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_users_timestamp();

----------------------------------------------------------------
-- CURRENCIES
----------------------------------------------------------------
DROP TABLE IF EXISTS currencies CASCADE;
CREATE TABLE currencies (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE, -- e.g. "US Dollar"
    symbol VARCHAR(5) NOT NULL UNIQUE -- e.g. "USD", "GHS"
);

-- seed 4 currencies: GHS, USD, EUR, GBP
INSERT INTO
    currencies (name, symbol)
VALUES ('Ghana Cedis', 'GHS'),
    ('US Dollar', 'USD'),
    ('Euro', 'EUR'),
    ('Pounds', 'GBP');

DROP FUNCTION IF EXISTS get_currency_info;
CREATE OR REPLACE FUNCTION get_currency_info(
    p_currency_id BIGINT
)
RETURNS TABLE (
    name VARCHAR,
    symbol VARCHAR
) AS $$
BEGIN
    RETURN QUERY
    SELECT c.name, c.symbol
    FROM currencies c
    WHERE id = p_currency_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Currency not found';
    END IF;
END;
$$ LANGUAGE plpgsql;

DROP FUNCTION IF EXISTS get_all_currencies;
CREATE OR REPLACE FUNCTION get_all_currencies()
RETURNS SETOF currencies
AS $$
BEGIN
    RETURN QUERY
    SELECT * FROM currencies;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'No currencies found';
    END IF;
END;
$$ LANGUAGE plpgsql;

----------------------------------------------------------------
-- CATEGORIES
----------------------------------------------------------------
