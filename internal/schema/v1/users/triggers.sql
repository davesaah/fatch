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
    NEW.updated_at := CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_users_update_timestamp ON users;
CREATE TRIGGER trg_users_update_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_users_timestamp();
