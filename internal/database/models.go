package database

import (
	"github.com/jackc/pgx/v5/pgtype"
)

// CreateUserParams contains the input parameters for the createUser function.
type CreateUserParams struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Passwd   string `json:"passwd"`
}

// ChangePasswordParams contains the input parameters for the changePassword function.
type ChangePasswordParams struct {
	UserID    pgtype.UUID `json:"-"`
	OldPasswd string      `json:"old_passwd"`
	NewPasswd string      `json:"new_passwd"`
}

// VerifyPasswordParams contains the input parameters for the verifyPassword function.
type VerifyPasswordParams struct {
	Username string `json:"username"`
	Passwd   string `json:"passwd"`
}

// GetUserByIdRow contains the output row for the getUserById function.
type GetUserByIdRow struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type GetCurrencyByIdRow struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type GetAllCurrenciesRow struct {
	CurrencyID int    `json:"currency_id"`
	Name       string `json:"name"`
	Symbol     string `json:"symbol"`
}

type CreateAccountParams struct {
	UserID      pgtype.UUID `json:"-"`
	AccountName string      `json:"account_name"`
	CurrencyID  int         `json:"currency_id"`
	Balance     float64     `json:"balance"`
	Description string      `json:"description"`
}

type GetAccountDetailsParams struct {
	UserID    pgtype.UUID `json:"-"`
	AccountID int         `json:"account_id"`
}

type GetAccountDetailsRow struct {
	AccountName string  `json:"account_name"`
	Balance     float64 `json:"balance"`
	Currency    string  `json:"currency"`
	Description string  `json:"description"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	IsArchived  bool    `json:"is_archived"`
}

type GetAllUserAccountsRow struct {
	AccountID int `json:"account_id"`
	GetAccountDetailsRow
}

type ArchiveAccountByIDParams struct {
	GetAccountDetailsParams
}
