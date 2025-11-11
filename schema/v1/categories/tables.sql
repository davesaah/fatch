CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    type category_type NOT NULL,
    name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE subcategories (
    id BIGSERIAL PRIMARY KEY,
    category_id BIGINT NOT NULL REFERENCES categories (id) ON DELETE CASCADE,
    name VARCHAR(150) NOT NULL,
    UNIQUE (category_id, name) -- ensures name is unique only within its category
);

ALTER TABLE categories ADD COLUMN user_id uuid DEFAULT NULL REFERENCES users (id) ON DELETE CASCADE;

CREATE INDEX idx_categories_user_id ON categories (user_id);
CREATE INDEX idx_subcategories_category_id ON subcategories (category_id);
