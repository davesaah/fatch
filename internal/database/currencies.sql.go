package database

import "context"

func (q *Queries) GetCurrencyByID(ctx context.Context, currencyID int) (*GetCurrencyByIdRow, error) {
	row := q.db.QueryRow(ctx, getCurrencyById, currencyID)
	var i GetCurrencyByIdRow
	err := row.Scan(&i.Name, &i.Symbol)
	return &i, err
}

func (q *Queries) GetAllCurrencies(ctx context.Context) ([]GetAllCurrenciesRow, error) {
	rows, err := q.db.Query(ctx, getAllCurrencies)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var currencies []GetAllCurrenciesRow
	for rows.Next() {
		var currency GetAllCurrenciesRow
		if err := rows.Scan(&currency.CurrencyID, &currency.Name, &currency.Symbol); err != nil {
			return nil, err
		}
		currencies = append(currencies, currency)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return currencies, nil
}
