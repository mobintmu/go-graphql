-- name: CreateProduct :one
INSERT INTO products (product_name, product_description, price, is_active)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetProduct :one
SELECT * FROM products WHERE id = $1;

-- name: ListProducts :many
SELECT * FROM products ORDER BY created_at DESC;

-- name: UpdateProduct :one
UPDATE products
SET product_name = $2, product_description = $3, price = $4, is_active = $5
WHERE id = $1
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;

-- name: ListProductsWithFilters :many
SELECT id, product_name, product_description, price, is_active, created_at
FROM products
WHERE
  (sqlc.narg('id')::int IS NULL OR id = sqlc.narg('id'))
  AND (sqlc.narg('product_name')::text IS NULL OR product_name ILIKE '%' || sqlc.narg('product_name') || '%')
  AND (sqlc.narg('min_price')::bigint IS NULL OR price >= sqlc.narg('min_price'))
  AND (sqlc.narg('max_price')::bigint IS NULL OR price <= sqlc.narg('max_price'))
  AND (sqlc.narg('is_active')::bool IS NULL OR is_active = sqlc.narg('is_active'))
  AND (sqlc.narg('product_description')::text IS NULL OR product_description ILIKE '%' || sqlc.narg('product_description') || '%')
ORDER BY id
LIMIT sqlc.narg('limit')
OFFSET sqlc.narg('offset');


-- name: CountProductsWithFilters :one
SELECT COUNT(*) as count
FROM products
WHERE
  (sqlc.narg('id')::int IS NULL OR id = sqlc.narg('id'))
  AND (sqlc.narg('product_name')::text IS NULL OR product_name ILIKE '%' || sqlc.narg('product_name') || '%')
  AND (sqlc.narg('min_price')::bigint IS NULL OR price >= sqlc.narg('min_price'))
  AND (sqlc.narg('max_price')::bigint IS NULL OR price <= sqlc.narg('max_price'))
  AND (sqlc.narg('is_active')::bool IS NULL OR is_active = sqlc.narg('is_active'))
  AND (sqlc.narg('product_description')::text IS NULL OR product_description ILIKE '%' || sqlc.narg('product_description') || '%');