-- Active: 1758835250687@@127.0.0.1@5432@local@fatch
DROP FUNCTION create_account;

CREATE OR REPLACE FUNCTION create_account(
    p_user_id UUID,
    p_name VARCHAR,
    p_currency_id BIGINT,
    p_balance DECIMAL(18,2) DEFAULT 0.00,
    p_description TEXT DEFAULT NULL
)
RETURNS TABLE (
    account_id BIGINT,
    account_name VARCHAR,
    currency_name VARCHAR,
    account_balance DECIMAL(18,2),
    account_description TEXT,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    is_archived BOOLEAN
)
LANGUAGE plpgsql
AS $$
DECLARE
    account_id BIGINT;
BEGIN
    -- check if account name exist for current user
    IF EXISTS (
        SELECT 1
        FROM accounts
        WHERE user_id = p_user_id AND lower(name) = lower(p_name)
    ) THEN
        RAISE EXCEPTION 'Account with name "%" already exists for this user', p_name
            USING ERRCODE = 'unique_violation';
    END IF;

    -- insert fresh account
    -- get the id of the latest insert
    INSERT INTO accounts(user_id, name, balance, description, currency_id)
    VALUES (p_user_id, p_name, p_balance, p_description, p_currency_id)
    RETURNING id INTO account_id;

    -- return account details
    RETURN QUERY SELECT
        acc.id,
        acc.name AS account_name,
        curr.name AS account_currency_name,
        acc.balance AS account_balance,
        acc.description AS account_description,
        acc.created_at,
        acc.updated_at,
        acc.is_archived
    FROM accounts acc
    JOIN currencies curr ON acc.currency_id = curr.id
    WHERE acc.id = account_id;
END;
$$;

DROP FUNCTION get_account_details;

CREATE
OR REPLACE FUNCTION get_account_details (
    p_account_id bigint,
    p_user_id uuid
) RETURNS TABLE (
    account_name VARCHAR,
    currency_name VARCHAR,
    account_balance DECIMAL(18,2),
    account_description TEXT,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    is_archived BOOLEAN
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY SELECT
    acc.name AS account_name,
    curr.name AS currency_name,
    acc.balance AS account_balance,
    acc.description AS account_description,
    acc.created_at,
    acc.updated_at,
    acc.is_archived
    FROM accounts acc
    JOIN currencies curr
    ON acc.currency_id = curr.id
    WHERE acc.id = p_account_id AND acc.user_id = p_user_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Account not found';
    END IF;
END;
$$;

CREATE OR REPLACE FUNCTION archive_account_by_id (
    p_account_id bigint,
    p_user_id uuid
) RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE accounts
    SET is_archived = TRUE,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = p_account_id AND user_id = p_user_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Account not found';
    END IF;
END;
$$;

DROP FUNCTION get_all_user_accounts;

CREATE OR REPLACE FUNCTION get_all_user_accounts (p_user_id uuid)
RETURNS TABLE (
    id bigint,
    name varchar,
    balance decimal,
    currency varchar,
    description text,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    is_archived BOOLEAN
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY SELECT
    acc.id, acc.name, acc.balance, curr.name as currency,
    acc.description, acc.created_at, acc.updated_at,
    acc.is_archived
    FROM accounts acc
    JOIN currencies curr
    ON acc.currency_id = curr.id
    WHERE acc.user_id = p_user_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'No accounts found';
    END IF;
END;
$$;