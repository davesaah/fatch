-- v1 of DB Schema for fatch
DROP SCHEMA IF EXISTS fatch CASCADE;
CREATE SCHEMA fatch;
SET search_path TO fatch;

-- for monitoring
CREATE EXTENSION pg_stat_statements;

-- For gen_random_uuid()
-- Password hashing.
CREATE EXTENSION pgcrypto;

-- For case-insensitive uniqueness.
CREATE EXTENSION citext;

----------------------------------------------------------------
-- USERS 
----------------------------------------------------------------
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    username citext NOT NULL UNIQUE,
    email citext NOT NULL UNIQUE,
    passwd TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_users_username ON users (username);
CREATE INDEX idx_users_email ON users (email);

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
    NEW.updated_at := now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_users_update_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_users_timestamp();

----------------------------------------------------------------
-- CURRENCIES
----------------------------------------------------------------
CREATE TABLE currencies (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE, -- e.g. "US Dollar"
    symbol VARCHAR(5) NOT NULL UNIQUE -- e.g. "USD", "GHS"
);

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
CREATE TYPE category_type AS ENUM ('expense', 'income');

CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    user_id uuid DEFAULT NULL REFERENCES users (id) ON DELETE CASCADE,
    type category_type NOT NULL,
    name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE subcategories (
    id BIGSERIAL PRIMARY KEY,
    category_id BIGINT NOT NULL REFERENCES categories (id) ON DELETE CASCADE,
    name VARCHAR(150) NOT NULL,
    UNIQUE (category_id, name) -- ensures name is unique only within its category
);

CREATE INDEX idx_categories_user_id ON categories (user_id);
CREATE INDEX idx_subcategories_category_id ON subcategories (category_id);

CREATE OR REPLACE FUNCTION add_category(
    p_type VARCHAR,
    p_name VARCHAR,
    p_user_id UUID = NULL
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
	-- Validate category type
    BEGIN
        PERFORM p_type::category_type;  -- Try to cast to ENUM
    EXCEPTION
        WHEN others THEN
            RAISE EXCEPTION 'Category type "%" does not exist', p_type;
    END;

	BEGIN
    	INSERT INTO categories(type, name, user_id)
    	VALUES (p_type::category_type, p_name, p_user_id);

	EXCEPTION
    	WHEN unique_violation THEN
        	RAISE EXCEPTION 'Category name "%" already exists', p_name;
	END;
END;
$$;

CREATE OR REPLACE FUNCTION add_subcategory(
    p_category_id BIGINT,
    p_name VARCHAR
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
	IF NOT EXISTS (
        SELECT 1 FROM categories WHERE id = p_category_id
    ) THEN
        RAISE EXCEPTION 'Category with ID % does not exist', p_category_id;
    END IF;

    INSERT INTO subcategories(category_id, name)
    VALUES (p_category_id, p_name);
EXCEPTION
    WHEN unique_violation THEN
        RAISE EXCEPTION 'Subcategory name "%" already exists in this category', p_name;
END;
$$;

CREATE OR REPLACE FUNCTION get_categories()
RETURNS TABLE (
    category_id BIGINT,
    category_name VARCHAR,
    category_type category_type
) AS $$
BEGIN
    RETURN QUERY
    SELECT c.id, c.name, c.type
    FROM categories c
    ORDER BY c.type, c.name;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_subcategories(
    p_category_id BIGINT
)
RETURNS TABLE (
    subcategory_id BIGINT,
    subcategory_name VARCHAR
) AS $$
BEGIN
    RETURN QUERY
    SELECT id, name
    FROM subcategories
    WHERE category_id = p_category_id
    ORDER BY name;

    IF NOT FOUND THEN
        RAISE NOTICE 'No subcategories found for category ID %', p_category_id;
    END IF;
END;
$$ LANGUAGE plpgsql;

----------------------------------------------------------------
-- ACCOUNTS
----------------------------------------------------------------
CREATE TABLE accounts (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name VARCHAR(25) NOT NULL,
    balance DECIMAL(18, 2) NOT NULL DEFAULT 0.00,
    description TEXT,
    currency_id BIGINT NOT NULL REFERENCES currencies (id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_archived BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_accounts_user_id_id ON accounts (user_id, id);
CREATE INDEX idx_accounts_user_id_name ON accounts (user_id, name);
CREATE INDEX idx_accounts_user_id ON accounts (user_id);

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

CREATE OR REPLACE FUNCTION get_account_details (
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

CREATE OR REPLACE FUNCTION update_accounts_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at := CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_accounts_update_timestamp
BEFORE UPDATE ON accounts
FOR EACH ROW
EXECUTE FUNCTION update_accounts_timestamp();

----------------------------------------------------------------
-- BUDGETS
----------------------------------------------------------------
CREATE TYPE budget_period AS ENUM ('day', 'week', 'month', 'year', 'once');
CREATE TYPE budget_target_type AS ENUM ('account', 'category', 'subcategory');

CREATE TABLE budgets (
    id BIGSERIAL PRIMARY KEY,
    budget_name VARCHAR NOT NULL,
    amount DECIMAL(18,2) NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    start_date DATE NOT NULL,
    frequency INT NOT NULL DEFAULT 1 CHECK (frequency > 0),  -- number of periods per cycle
    period budget_period NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE budget_rules (
    id BIGSERIAL PRIMARY KEY,
    budget_id BIGINT NOT NULL REFERENCES budgets(id) ON DELETE CASCADE,
    target_type budget_target_type NOT NULL,
    target_id BIGINT NOT NULL
);

CREATE INDEX idx_budget_rules_target ON budget_rules (target_type, target_id);
CREATE INDEX idx_budgets_user_id ON budgets (user_id);

CREATE OR REPLACE FUNCTION create_budget(
    p_user_id UUID,
    p_name VARCHAR,
    p_amount DECIMAL(18,2),
    p_start_date DATE,
    p_frequency INT DEFAULT 1,
    p_period budget_period DEFAULT 'month'
)
RETURNS VOID
AS $$
BEGIN
    INSERT INTO budgets(user_id, budget_name, amount, start_date, frequency, period)
    VALUES (p_user_id, p_name, p_amount, p_start_date, p_frequency, p_period);
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION add_budget_rule(
    p_budget_id BIGINT,
    p_target_type budget_target_type,
    p_target_id BIGINT
)
RETURNS VOID
AS $$
BEGIN
    INSERT INTO budget_rules(budget_id, target_type, target_id)
    VALUES (p_budget_id, p_target_type, p_target_id);
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_budgets_for_target(
   p_user_id UUID,
   p_target_type budget_target_type
)
RETURNS TABLE (
    budget_id BIGINT,
    budget_name VARCHAR,
    target_id BIGINT,
    amount DECIMAL(18,2),
    start_date DATE,
    frequency INT,
    period budget_period
) AS $$
BEGIN
    RETURN QUERY
    SELECT b.id, b.budget_name, r.target_id, b.amount, b.start_date, b.frequency, b.period
    FROM budgets b
    JOIN budget_rules r ON r.budget_id = b.id
    WHERE b.user_id = p_user_id
      AND r.target_type = p_target_type
    ORDER BY b.start_date DESC;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_budget_rules(
    p_user_id UUID,
    p_budget_id BIGINT
)
RETURNS TABLE (
    rule_id BIGINT,
    target_type budget_target_type,
    target_id BIGINT,
    amount DECIMAL(18,2),
    period budget_period,
    frequency INT,
    start_date DATE
) AS $$
BEGIN
    RETURN QUERY
    SELECT r.id, r.target_type, r.target_id, b.amount, b.period, b.frequency, b.start_date
    FROM budget_rules r
    JOIN budgets b ON b.id = r.budget_id
    WHERE b.user_id = p_user_id
      AND b.id = p_budget_id
    ORDER BY r.target_type, r.target_id;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_user_budgets(p_user_id BIGINT)
RETURNS TABLE (
    budget_id BIGINT,
    amount DECIMAL(18,2),
    budget_name VARCHAR,
    start_date DATE,
    frequency INT,
    period budget_period,
    created_at TIMESTAMPTZ
) AS $$
BEGIN
    RETURN QUERY
    SELECT b.id, b.amount, b.budget_name, b.start_date, b.frequency, b.period, b.created_at
    FROM budgets b
    WHERE user_id = p_user_id
    ORDER BY start_date DESC;
END;
$$ LANGUAGE plpgsql;

-- validation from a trigger
CREATE OR REPLACE FUNCTION validate_budget_rule_target()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.target_type = 'account' THEN
        IF NOT EXISTS (SELECT 1 FROM accounts WHERE id = NEW.target_id) THEN
            RAISE EXCEPTION 'Invalid target_id % for target_type account', NEW.target_id;
        END IF;
    ELSIF NEW.target_type = 'category' THEN
        IF NOT EXISTS (SELECT 1 FROM categories WHERE id = NEW.target_id) THEN
            RAISE EXCEPTION 'Invalid target_id % for target_type category', NEW.target_id;
        END IF;
    ELSIF NEW.target_type = 'subcategory' THEN
        IF NOT EXISTS (SELECT 1 FROM subcategories WHERE id = NEW.target_id) THEN
            RAISE EXCEPTION 'Invalid target_id % for target_type subcategory', NEW.target_id;
        END IF;
    ELSE
        RAISE EXCEPTION 'Invalid target_type %', NEW.target_type;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- attach trigger to budget rules
CREATE TRIGGER trg_validate_budget_rule_target
BEFORE INSERT OR UPDATE ON budget_rules
FOR EACH ROW
EXECUTE FUNCTION validate_budget_rule_target();
