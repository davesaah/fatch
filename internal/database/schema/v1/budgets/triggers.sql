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
