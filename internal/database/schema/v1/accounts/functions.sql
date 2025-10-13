-- Active: 1758835250687@@127.0.0.1@5432@local@fatch
CREATE OR REPLACE FUNCTION get_account_balance(
    p_account_id BIGINT,
    p_user_id UUID
)
RETURNS DECIMAL(18,2) AS $$
DECLARE
    _balance DECIMAL(18,2);
BEGIN
    SELECT balance INTO _balance
    FROM accounts
    WHERE id = p_account_id
      AND user_id = p_user_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Account not found for the given user';
    END IF;

    RETURN _balance;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION create_account(
    p_user_id UUID,
    p_name VARCHAR,
    p_currency_id BIGINT,
    p_balance DECIMAL(18,2) DEFAULT 0.00,
    p_description TEXT DEFAULT NULL
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    -- check if account name exist for current user
    IF EXISTS (
        SELECT 1
        FROM accounts
        WHERE user_id = p_user_id AND name = p_name
    ) THEN
        RAISE EXCEPTION 'Account with name "%" already exists for this user', p_name
            USING ERRCODE = 'unique_violation';
    END IF;

    -- insert fresh account
    INSERT INTO accounts(user_id, name, balance, description, currency_id)
    VALUES (p_user_id, p_name, p_balance, p_description, p_currency_id);
END;
$$;

CREATE
OR REPLACE FUNCTION get_account_details (
    p_account_id bigint,
    p_user_id uuid
) RETURNS TABLE (
    name varchar,
    balance decimal,
    currency varchar,
    description text,
    created_at timestamptz,
    updated_at timestamptz
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY SELECT
    acc.name, acc.balance, curr.name as currency,
    acc.description, acc.created_at, acc.updated_at
    FROM accounts acc
    JOIN currencies curr
    ON acc.currency_id = curr.id
    WHERE acc.id = p_account_id AND acc.user_id = p_user_id;
END;
$$;

CREATE OR REPLACE FUNCTION get_all_user_accounts (p_user_id uuid)
RETURNS TABLE (
    id bigint,
    name varchar,
    balance decimal,
    currency varchar,
    description text,
    created_at timestamptz,
    updated_at timestamptz
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY SELECT
    acc.id, acc.name, acc.balance, curr.name as currency,
    acc.description, acc.created_at, acc.updated_at
    FROM accounts acc
    JOIN currencies curr
    ON acc.currency_id = curr.id
    WHERE acc.user_id = p_user_id;
END;
$$;