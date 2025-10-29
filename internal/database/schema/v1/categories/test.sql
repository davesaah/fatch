select add_category('budget', 'housing');
select add_category('expense', 'housing');
select add_category('expense', 'housing');

select add_subcategory(1, 'Utilities');
select add_subcategory(2, 'Utilities');

select get_categories();
select get_subcategories(1);
