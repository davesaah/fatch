CREATE TABLE budgets (
    id BIGSERIAL PRIMARY KEY,
    budget_name VARCHAR NOT NULL,
    amount DECIMAL(18,2) NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    start_date DATE NOT NULL,
    frequency INT NOT NULL DEFAULT 1 CHECK (frequency > 0),       -- number of periods per cycle
    period budget_period NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE budget_rules (
    id BIGSERIAL PRIMARY KEY,
    budget_id BIGINT NOT NULL REFERENCES budgets(id) ON DELETE CASCADE,
    target_type budget_target_type NOT NULL,
    target_id BIGINT NOT NULL
);

-- Index for fast lookups by target_type and target_id
CREATE INDEX idx_budget_rules_target ON budget_rules (target_type, target_id);
CREATE INDEX idx_budget_user_id ON budget_rules (user_id);
