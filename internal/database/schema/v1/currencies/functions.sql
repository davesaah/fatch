-- Active: 1758835250687@@127.0.0.1@5432@local@fatch
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
