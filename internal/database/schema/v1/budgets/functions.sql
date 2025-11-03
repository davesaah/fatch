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

-- CREATE OR REPLACE FUNCTION get_budget_consumption(p_budget_id BIGINT)
-- RETURNS DECIMAL(18,2) AS $$
-- DECLARE
--     _budget budgets%ROWTYPE;
--     _period_start DATE;
--     _period_end DATE;
--     _now DATE := CURRENT_DATE;
--     _total DECIMAL(18,2);
-- BEGIN
--     -- Get budget details
--     SELECT * INTO _budget
--     FROM budgets
--     WHERE id = _budget_id;

--     IF NOT FOUND THEN
--         RAISE EXCEPTION 'Budget % not found', _budget_id;
--     END IF;

--     -- Determine current period start and end
--     IF _budget.period = 'once' THEN
--         _period_start := _budget.start_date;
--         _period_end := _budget.start_date;
--     ELSE
--         CASE _budget.period
--             WHEN 'day' THEN
--                 _period_start := _budget.start_date + ((_now - _budget.start_date)::int / _budget.frequency) * _budget.frequency;
--                 _period_end := _period_start + _budget.frequency - 1;
--             WHEN 'week' THEN
--                 _period_start := _budget.start_date + ((EXTRACT(EPOCH FROM _now - _budget.start_date) / 86400)::int / (7 * _budget.frequency)) * 7 * _budget.frequency;
--                 _period_end := _period_start + (7 * _budget.frequency) - 1;
--             WHEN 'month' THEN
--                 _period_start := (DATE_TRUNC('month', _budget.start_date) + INTERVAL '1 month' * FLOOR((DATE_PART('year', _now) - DATE_PART('year', _budget.start_date)) * 12 + (DATE_PART('month', _now) - DATE_PART('month', _budget.start_date))) / _budget.frequency)::date;
--                 _period_end := (_period_start + INTERVAL '1 month' * _budget.frequency - INTERVAL '1 day')::date;
--             WHEN 'year' THEN
--                 _period_start := (DATE_TRUNC('year', _budget.start_date) + INTERVAL '1 year' * FLOOR((DATE_PART('year', _now) - DATE_PART('year', _budget.start_date)) / _budget.frequency))::date;
--                 _period_end := (_period_start + INTERVAL '1 year' * _budget.frequency - INTERVAL '1 day')::date;
--         END CASE;
--     END IF;

--     -- Sum all transactions that match this budget in one query
--     SELECT COALESCE(SUM(t.amount), 0) INTO _total
--     FROM transactions t
--     JOIN budget_rules r ON
--         (r.budget_id = _budget.id AND
--          ((r.target_type = 'account' AND (t.debit_account_id = r.target_id OR t.credit_account_id = r.target_id)) OR
--           (r.target_type = 'category' AND t.category_id = r.target_id) OR
--           (r.target_type = 'subcategory' AND t.subcategory_id = r.target_id)))
--     WHERE t.created_at::date BETWEEN _period_start AND _period_end
--       AND (t.type = 'expense' OR t.type = 'transfer');

--     RETURN _total;
-- END;
-- $$ LANGUAGE plpgsql;

-- CREATE OR REPLACE FUNCTION get_budget_report(p_budget_id BIGINT)
-- RETURNS TABLE (
--     budget_id BIGINT,
--     user_id BIGINT,
--     total_amount DECIMAL(18,2),
--     spent_amount DECIMAL(18,2),
--     remaining_amount DECIMAL(18,2),
--     percent_consumed NUMERIC(5,2),
--     period_start DATE,
--     period_end DATE
-- ) AS $$
-- DECLARE
--     _budget budgets%ROWTYPE;
--     _period_start DATE;
--     _period_end DATE;
--     _now DATE := CURRENT_DATE;
--     _spent DECIMAL(18,2);
-- BEGIN
--     -- Get budget details
--     SELECT * INTO _budget
--     FROM budgets
--     WHERE id = p_budget_id;

--     IF NOT FOUND THEN
--         RAISE EXCEPTION 'Budget % not found', p_budget_id;
--     END IF;

--     -- Determine current period
--     IF _budget.period = 'once' THEN
--         _period_start := _budget.start_date;
--         _period_end := _budget.start_date;
--     ELSE
--         CASE _budget.period
--             WHEN 'day' THEN
--                 _period_start := _budget.start_date + ((_now - _budget.start_date)::int / _budget.frequency) * _budget.frequency;
--                 _period_end := _period_start + _budget.frequency - 1;
--             WHEN 'week' THEN
--                 _period_start := _budget.start_date + ((EXTRACT(EPOCH FROM _now - _budget.start_date) / 86400)::int / (7 * _budget.frequency)) * 7 * _budget.frequency;
--                 _period_end := _period_start + (7 * _budget.frequency) - 1;
--             WHEN 'month' THEN
--                 _period_start := (DATE_TRUNC('month', _budget.start_date) + INTERVAL '1 month' * FLOOR((DATE_PART('year', _now) - DATE_PART('year', _budget.start_date)) * 12 + (DATE_PART('month', _now) - DATE_PART('month', _budget.start_date))) / _budget.frequency)::date;
--                 _period_end := (_period_start + INTERVAL '1 month' * _budget.frequency - INTERVAL '1 day')::date;
--             WHEN 'year' THEN
--                 _period_start := (DATE_TRUNC('year', _budget.start_date) + INTERVAL '1 year' * FLOOR((DATE_PART('year', v_now) - DATE_PART('year', _budget.start_date)) / _budget.frequency))::date;
--                 _period_end := (_period_start + INTERVAL '1 year' * _budget.frequency - INTERVAL '1 day')::date;
--         END CASE;
--     END IF;

--     -- Sum all matching transactions in the period
--     SELECT COALESCE(SUM(t.amount), 0) INTO _spent
--     FROM transactions t
--     JOIN budget_rules r ON
--         r.budget_id = _budget.id AND
--         ((r.target_type = 'account' AND (t.debit_account_id = r.target_id OR t.credit_account_id = r.target_id)) OR
--          (r.target_type = 'category' AND t.category_id = r.target_id) OR
--          (r.target_type = 'subcategory' AND t.subcategory_id = r.target_id))
--     WHERE t.created_at::date BETWEEN _period_start AND _period_end
--       AND (t.type = 'expense' OR t.type = 'transfer');

--     -- Return all relevant info
--     RETURN QUERY
--     SELECT
--         _budget.id AS budget_id,
--         _budget.user_id,
--         _budget.amount AS total_amount,
--         _spent AS spent_amount,
--         _budget.amount - v_spent AS remaining_amount,
--         CASE WHEN _budget.amount = 0 THEN 0 ELSE ROUND((_spent / _budget.amount) * 100, 2) END AS percent_consumed,
--         _period_start,
--         _period_end;
-- END;
-- $$ LANGUAGE plpgsql;

-- CREATE OR REPLACE FUNCTION get_user_budget_report(p_user_id BIGINT)
-- RETURNS TABLE (
--     budget_id BIGINT,
--     total_amount DECIMAL(18,2),
--     spent_amount DECIMAL(18,2),
--     remaining_amount DECIMAL(18,2),
--     percent_consumed NUMERIC(5,2),
--     period_start DATE,
--     period_end DATE
-- ) AS $$
-- BEGIN
--     RETURN QUERY
--     SELECT gr.*
--     FROM budgets b
--     CROSS JOIN LATERAL get_budget_report(b.id) gr
--     WHERE b.user_id = p_user_id;
-- END;
-- $$ LANGUAGE plpgsql;
