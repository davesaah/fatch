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
