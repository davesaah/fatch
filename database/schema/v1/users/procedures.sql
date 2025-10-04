-- Active: 1758835250687@@127.0.0.1@5432@local@fatch
DROP PROCEDURE IF EXISTS create_user (citext, citext, TEXT);

CREATE OR REPLACE PROCEDURE create_user(
    p_username citext,
    p_email citext,
    p_passwd TEXT
) AS $$
BEGIN
    INSERT INTO users(username, email, passwd)
    VALUES (p_username, p_email, p_passwd);
EXCEPTION
    WHEN unique_violation THEN
        RAISE EXCEPTION 'Username or email already exists';
END;
$$ LANGUAGE plpgsql;

CREATE
OR REPLACE PROCEDURE change_password (
    p_user_id UUID,
    p_old_passwd TEXT,
    p_new_passwd TEXT
) AS $$
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