package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getCategories = "SELECT (get_categories($1::uuid)).*"

func (q *Queries) GetCategories(ctx context.Context, userId pgtype.UUID) ([]GetAllCategoriesRow, error) {
	rows, err := q.db.Query(ctx, getCategories, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []GetAllCategoriesRow
	for rows.Next() {
		var category GetAllCategoriesRow
		if err := rows.Scan(&category.CategoryID, &category.CategoryName, &category.Type); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

const getCategoryByID = "SELECT (get_category_by_id($1::uuid, $2::bigint)).*"

func (q *Queries) GetCategoryByID(ctx context.Context, arg GetCategoryByIDParams) (*GetCategoryByIDRow, error) {
	row := q.db.QueryRow(ctx, getCategoryByID, arg.UserID, arg.CategoryID)
	var category GetCategoryByIDRow
	if err := row.Scan(&category.CategoryName, &category.Type); err != nil {
		return nil, err
	}
	return &category, nil
}

const addCategory = "SELECT add_category($1::varchar, $2::varchar, $3::uuid)"

func (q *Queries) AddCategory(ctx context.Context, arg CreateCategoryParams) error {
	_, err := q.db.Exec(ctx, addCategory, arg.Type, arg.CategoryName, arg.UserID)
	return err
}

const updateCategory = "SELECT update_category($1::bigint, $2::varchar, $3::varchar, $4::uuid)"

func (q *Queries) UpdateCategory(ctx context.Context, arg UpdateCategoryParams) error {
	_, err := q.db.Exec(ctx, updateCategory, arg.CategoryID, arg.Type, arg.CategoryName, arg.UserID)
	return err
}

const deleteCategory = "SELECT delete_category($1::bigint, $2::uuid)"

func (q *Queries) DeleteCategory(ctx context.Context, arg DeleteCategoryParams) error {
	_, err := q.db.Exec(ctx, deleteCategory, arg.CategoryID, arg.UserID)
	return err
}
