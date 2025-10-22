package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

// CreateUser creates a new user in the database.
func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.Exec(ctx, createUser, arg.Username, arg.Email, arg.Passwd)
	return err
}

// ChangePassword changes the password for a user in the database.
func (q *Queries) ChangePassword(ctx context.Context, arg ChangePasswordParams) error {
	_, err := q.db.Exec(ctx, changePassword, arg.UserID, arg.OldPasswd, arg.NewPasswd)
	return err
}

// VerifyPassword verifies a user's password.
func (q *Queries) VerifyPassword(ctx context.Context, arg VerifyPasswordParams) (*pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, verifyPassword, arg.Username, arg.Passwd)
	var i pgtype.UUID
	err := row.Scan(&i)
	return &i, err
}

// GetUserById retrieves a user by their ID.
func (q *Queries) GetUserById(ctx context.Context, userID pgtype.UUID) (*GetUserByIdRow, error) {
	row := q.db.QueryRow(ctx, getUserById, userID)
	var i GetUserByIdRow
	err := row.Scan(&i.Username, &i.Email)
	return &i, err
}
