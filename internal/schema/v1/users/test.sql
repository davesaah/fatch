-- Active: 1758835250687@@127.0.0.1@5432@local@fatch
-- Flow 1: Create a user

SELECT fatch.create_user (
        'someone', 'someone@gmail.com', 'password123'
    );

SELECT fatch.create_user (
        'davesaah', 'davesaah@gmail.com', 'password123'
    );

SELECT fatch.create_user (
        'davesaah', 'davesaah@gmail.com', 'password12345'
    );
-- Flow 2: Verify password
SELECT fatch.verify_password ('davesaah', '', 'password123');

SELECT fatch.verify_password (
        '', 'davesaah@gmail.com', 'password123'
    );

SELECT fatch.verify_password ('someone', '', 'password');

SELECT fatch.verify_password ('davesaah', '', 'pasword123');

SELECT fatch.verify_password (
        '', 'davesaah@gmail.com', 'pasword123'
    );

SELECT fatch.verify_password ( 'davesaah', '', 'newpassword456' );

-- Flow 3: Change password
SELECT fatch.change_password (
        '360ef0fa-75a0-4fc3-b842-00fb4725b353', 'passwod123', 'newpassword456'
    );

SELECT fatch.change_password (
        '330ef0fa-75a0-4fc3-b842-00fb4725b353', 'passwod123', 'newpassword456'
    );

SELECT fatch.change_password (
        '360ef0fa-75a0-4fc3-b842-00fb4725b353', 'password123', 'newpassword456'
    );

SELECT get_user_by_id ( '487e1477-d2f6-4d22-bf16-b87cb4b4b9df' );