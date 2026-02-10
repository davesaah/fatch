package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = "SELECT create_user ($1::citext, $2::citext, $3::text)"

// CreateUser creates a new user in the database and creates an OTP
func (q *Queries) CreateUser(ctx context.Context, arg RegisterUserParams) (int, error) {
	row := q.db.QueryRow(ctx, createUser, arg.Username, arg.Email, arg.Passwd)
	var otp int
	err := row.Scan(&otp)
	return otp, err
}

const changePassword = "SELECT change_password ($1::uuid, $2::text, $3::text)"

// ChangePassword changes the password for a user in the database.
func (q *Queries) ChangePassword(ctx context.Context, arg ChangePasswordParams) error {
	_, err := q.db.Exec(ctx, changePassword, arg.UserID, arg.OldPasswd, arg.NewPasswd)
	return err
}

const verifyPassword = "SELECT verify_password ($1::citext, $2::text)"

// VerifyPassword verifies a user's password.
func (q *Queries) VerifyPassword(ctx context.Context, arg LoginParams) (*pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, verifyPassword, arg.Username, arg.Passwd)
	var i pgtype.UUID
	err := row.Scan(&i)
	return &i, err
}

const getUserByID = "SELECT (get_user_by_id($1::uuid)).*"

// GetUserByID retrieves a user by their ID.
func (q *Queries) GetUserByID(ctx context.Context, userID pgtype.UUID) (*GetUserByIDRow, error) {
	row := q.db.QueryRow(ctx, getUserByID, userID)
	var i GetUserByIDRow
	err := row.Scan(&i.Username, &i.Email)
	return &i, err
}

const verifyUser = "SELECT verify_user ($1::citext, $2::text, $3::int)"

// VerifyUser verifies a user's OTP.
func (q *Queries) VerifyUser(ctx context.Context, arg VerifyUserParams) error {
	_, err := q.db.Exec(ctx, verifyUser, arg.Username, arg.Passwd, arg.OTP)
	return err
}

const deleteUser = "SELECT delete_user ($1::uuid, $2::text)"

// DeleteUser deletes a user from the database.
func (q *Queries) DeleteUser(ctx context.Context, arg DeleteUserParams) error {
	_, err := q.db.Exec(ctx, deleteUser, arg.UserID, arg.Passwd)
	return err
}
