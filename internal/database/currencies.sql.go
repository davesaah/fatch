package database

import "context"

const getCurrencyByID = "SELECT (get_currency_info($1::bigint)).*"

func (q *Queries) GetCurrencyByID(
	ctx context.Context, currencyID int,
) (*GetCurrencyByIDRow, error) {
	row := q.db.QueryRow(ctx, getCurrencyByID, currencyID)
	var i GetCurrencyByIDRow
	err := row.Scan(&i.Name, &i.Symbol)
	return &i, err
}

const getAllCurrencies = "SELECT (get_all_currencies()).*"

func (q *Queries) GetAllCurrencies(ctx context.Context) ([]GetAllCurrenciesRow, error) {
	rows, err := q.db.Query(ctx, getAllCurrencies)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var currencies []GetAllCurrenciesRow
	for rows.Next() {
		var currency GetAllCurrenciesRow
		err := rows.Scan(&currency.CurrencyID, &currency.Name, &currency.Symbol)
		if err != nil {
			return nil, err
		}
		currencies = append(currencies, currency)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return currencies, nil
}
