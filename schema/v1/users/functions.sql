-- Active: 1758835250687@@127.0.0.1@5432@local@fatch

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

CREATE
OR REPLACE FUNCTION change_password (
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
