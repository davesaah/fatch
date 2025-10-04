-- Active: 1758835250687@@127.0.0.1@5432@local@fatch
DROP FUNCTION IF EXISTS verify_password (citext, citext, TEXT);

CREATE OR REPLACE FUNCTION verify_password(
    p_username citext,
    p_email citext,
    p_passwd TEXT
) RETURNS TABLE (
    is_valid BOOLEAN,
    user_id UUID
)
AS $$
DECLARE
    stored_hash TEXT;
    user_id UUID;
    is_valid BOOLEAN;
BEGIN
    SELECT passwd, id INTO stored_hash, user_id
    FROM users
    WHERE username = p_username OR email = p_email
    LIMIT 1;

    IF stored_hash IS NULL THEN
        is_valid := FALSE;
        user_id := NULL;
    ELSE
        is_valid := stored_hash = crypt(p_passwd, stored_hash);
    END IF;

    IF NOT is_valid THEN
        RETURN QUERY SELECT FALSE, NULL::UUID;
    ELSE
        RETURN QUERY SELECT TRUE, user_id;
    END IF;
END;
$$ LANGUAGE plpgsql;