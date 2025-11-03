select create_budget('77b0f090-2d4d-4d20-a034-6b1fce56e503', 20, CURRENT_DATE);
select add_budget_rule(1, 'account', 1);
select get_budgets_for_target('77b0f090-2d4d-4d20-a034-6b1fce56e503','account',1);
select get_budget_rules('77b0f090-2d4d-4d20-a034-6b1fce56e503', 1);
