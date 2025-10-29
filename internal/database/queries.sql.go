package database

// USERS
const createUser = "SELECT create_user ($1::citext, $2::citext, $3::text)"
const changePassword = "SELECT change_password ($1::uuid, $2::text, $3::text)"
const verifyPassword = "SELECT verify_password ($1::citext, $2::text)"
const getUserById = "SELECT (get_user_by_id($1::uuid)).*"

// CURRENCIES
const getCurrencyById = "SELECT (get_currency_info($1::bigint)).*"
const getAllCurrencies = "SELECT (get_all_currencies()).*"

// ACCOUNTS
const createAccount = "SELECT (create_account ($1::uuid, $2::varchar, $3::bigint, $4::decimal, $5::text)).*"
const getAccountDetails = "SELECT (get_account_details($1::bigint, $2::uuid)).*"
const getAllUserAccounts = "SELECT (get_all_user_accounts($1::uuid)).*"
const archiveAccountByID = "SELECT archive_account_by_id($1::bigint, $2::uuid)"

// CATEGORIES
const getCategories = "SELECT (get_categories()).*"
const getSubcategories = "SELECT (get_subcategories($1::bigint)).*"
const addCategory = "SELECT add_category($1::varchar, $2::varchar, $3::uuid)"
const addSubcategory = "SELECT add_subcategory($1::uuid, $2::varchar, $3::bigint)"
