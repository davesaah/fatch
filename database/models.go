package database

import (
	"github.com/jackc/pgx/v5/pgtype"
)

// USERS

type CreateUserParams struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Passwd   string `json:"passwd"`
}

type ChangePasswordParams struct {
	UserID    pgtype.UUID `json:"-"`
	OldPasswd string      `json:"old_passwd"`
	NewPasswd string      `json:"new_passwd"`
}

type VerifyPasswordParams struct {
	Username string `json:"username"`
	Passwd   string `json:"passwd"`
}

type GetUserByIDRow struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// CURRENCY

type GetCurrencyByIDRow struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type GetAllCurrenciesRow struct {
	CurrencyID int    `json:"currency_id"`
	Name       string `json:"name"`
	Symbol     string `json:"symbol"`
}

// ACCOUNTS

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
	IsArchive bool `json:"is_archive"`
}
