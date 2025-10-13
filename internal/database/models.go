package database

import "github.com/jackc/pgx/v5/pgtype"

// CreateUserParams contains the input parameters for the createUser function.
type CreateUserParams struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Passwd   string `json:"passwd"`
}

// ChangePasswordParams contains the input parameters for the changePassword function.
type ChangePasswordParams struct {
	UserID    pgtype.UUID `json:"user_id"`
	OldPasswd string      `json:"old_passwd"`
	NewPasswd string      `json:"new_passwd"`
}

// VerifyPasswordParams contains the input parameters for the verifyPassword function.
type VerifyPasswordParams struct {
	Username string `json:"username"`
	Passwd   string `json:"passwd"`
}

// VerifyPasswordRow contains the output row for the verifyPassword function.
type VerifyPasswordRow struct {
	UserID pgtype.UUID `json:"user_id"`
}

// GetUserByIdRow contains the output row for the getUserById function.
type GetUserByIdRow struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
