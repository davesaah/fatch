package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// TODO: update this guy to just return the account id
const createAccount = "SELECT (create_account ($1::uuid, $2::varchar, $3::bigint, $4::decimal, $5::text)).*"

func (q *Queries) CreateAccount(
	ctx context.Context, arg CreateAccountParams,
) (*GetAllUserAccountsRow, error) {
	row := q.db.QueryRow(
		ctx, createAccount, arg.UserID, arg.AccountName,
		arg.CurrencyID, arg.Balance, arg.Description,
	)
	var i GetAllUserAccountsRow
	var createdAt time.Time
	var updatedAt time.Time

	err := row.Scan(
		&i.AccountID, &i.AccountName, &i.Currency, &i.Balance,
		&i.Description, &createdAt, &updatedAt, &i.IsArchived,
	)

	i.CreatedAt = createdAt.Format(time.DateTime)
	i.UpdatedAt = updatedAt.Format(time.DateTime)

	return &i, err
}

const getAccountDetails = "SELECT (get_account_details($1::bigint, $2::uuid)).*"

func (q *Queries) GetAccountDetails(
	ctx context.Context, arg GetAccountDetailsParams,
) (*GetAccountDetailsRow, error) {
	row := q.db.QueryRow(ctx, getAccountDetails, arg.AccountID, arg.UserID)
	var i GetAccountDetailsRow

	var createdAt time.Time
	var updatedAt time.Time

	err := row.Scan(
		&i.AccountName, &i.Currency, &i.Balance,
		&i.Description, &createdAt, &updatedAt, &i.IsArchived,
	)

	i.CreatedAt = createdAt.Format(time.DateTime)
	i.UpdatedAt = updatedAt.Format(time.DateTime)

	return &i, err
}

const getAllUserAccounts = "SELECT (get_all_user_accounts($1::uuid)).*"

func (q *Queries) GetAllUserAccounts(
	ctx context.Context, userID pgtype.UUID,
) ([]GetAllUserAccountsRow, error) {
	rows, err := q.db.Query(ctx, getAllUserAccounts, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []GetAllUserAccountsRow
	var createdAt time.Time
	var updatedAt time.Time

	for rows.Next() {
		var account GetAllUserAccountsRow

		if err := rows.Scan(
			&account.AccountID,
			&account.AccountName,
			&account.Balance,
			&account.Currency,
			&account.Description,
			&createdAt,
			&updatedAt,
			&account.IsArchived); err != nil {
			return nil, err
		}

		account.CreatedAt = createdAt.Format(time.DateTime)
		account.UpdatedAt = updatedAt.Format(time.DateTime)
		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

const archiveAccountByID = "SELECT archive_account_by_id($1::bigint, $2::uuid, $3::boolean)"

func (q *Queries) ArchiveAccountByID(ctx context.Context, arg ArchiveAccountByIDParams) error {
	_, err := q.db.Exec(ctx, archiveAccountByID, arg.AccountID, arg.UserID, arg.IsArchive)
	return err
}
