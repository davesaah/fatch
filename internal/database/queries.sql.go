package database

// TODO: Delete file after transfer of queries

// CATEGORIES
const (
	getCategories    = "SELECT (get_categories()).*"
	getSubcategories = "SELECT (get_subcategories($1::bigint)).*"
	addCategory      = "SELECT add_category($1::varchar, $2::varchar, $3::uuid)"
	addSubcategory   = "SELECT add_subcategory($1::uuid, $2::varchar, $3::bigint)"
)
