-- Active: 1758835250687@@127.0.0.1@5432@local@fatch
-- Trigger function to hash passwords using bcrypt
DROP TRIGGER IF EXISTS trg_hash_password ON users;

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
CREATE TRIGGER trg_hash_password
BEFORE INSERT OR UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION hash_user_password();

-- Trigger function to update updated_at
CREATE OR REPLACE FUNCTION update_users_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at := CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_users_update_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_users_timestamp();