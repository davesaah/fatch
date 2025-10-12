package database

// USERS
const createUser = "SELECT create_user ($1::citext, $2::citext, $3::text)"
const changePassword = "SELECT change_password ($1::uuid, $2::text, $3::text)"
const verifyPassword = "SELECT verify_password ($1::citext, $2::text)"
const getUserById = "SELECT (get_user_by_id($1::uuid)).*"
