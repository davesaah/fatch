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

-- index for faster loopups
CREATE INDEX idx_accounts_user_id_id ON accounts (user_id, id);

CREATE INDEX idx_accounts_user_id_name ON accounts (user_id, name);

CREATE INDEX idx_accounts_user_id ON accounts (user_id);
