SET search_path TO fatch;

----------------------------------------------------------------
-- 4 currencies: GHS, USD, EUR, GBP
----------------------------------------------------------------
INSERT INTO
    currencies (name, symbol)
VALUES ('Ghana Cedis', 'GHS'),
    ('US Dollar', 'USD'),
    ('Euro', 'EUR'),
    ('Pounds', 'GBP');

----------------------------------------------------------------
-- CATEGORIES
----------------------------------------------------------------
select add_category('income', 'Gifts');
select add_category('income', 'Salary');
select add_category('income', 'Commission');
select add_category('income', 'ROI');
select add_category('income', 'Sales');
select add_category('income', 'Contract');

select add_category('expense', 'Shopping');
select add_category('expense', 'Food & Drinks');
select add_category('expense', 'Housing');
select add_category('expense', 'Bills');
select add_category('expense', 'Transport');
select add_category('expense', 'Lifestyle');
select add_category('expense', 'God''s Projects');
select add_category('expense', 'Processing Error');
select add_category('expense', 'Stocks');

----------------------------------------------------------------
-- SUBCATEGORIES
----------------------------------------------------------------
-- select add_subcategory(1, '');
