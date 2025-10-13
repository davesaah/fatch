-- Active: 1758835250687@@127.0.0.1@5432@local@fatch
SELECT verify_password ('lael', '2xkWa24ABN');

SELECT get_account_balance (
        1, '77b0f090-2d4d-4d20-a034-6b1fce56e503'
    );

SELECT create_account (
        '77b0f090-2d4d-4d20-a034-6b1fce56e503', 'ecobank', 1
    );

SELECT get_account_details (
        1, '77b0f090-2d4d-4d20-a034-6b1fce56e503'
    );

SELECT get_all_user_accounts (
        '77b0f090-2d4d-4d20-a034-6b1fce56e503'
    );