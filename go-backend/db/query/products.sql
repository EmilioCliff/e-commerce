-- CreateProduct :one
INSERT INTO products (
    product_name, description, price, quantity, discount, rating, size_options, color_options, category, brand, image_url
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
RETURNING *;

-- name: GetAllProducts :many
SELECT * FROM products
ORDER BY product_name;

-- name: GetProductByCategory :many
SELECT * FROM products
WHERE category = $1;

-- name: GetProductByUpdatedTime :many
SELECT * FROM products
ORDER BY updated_at DESC;

-- GetProduct :one
SELECT * FROM products
WHERE id = $1
LIMIT 1;

-- GetProductForUpdate: one
SELECT  * FROM products
WHERE id = $1
FOR NO KEY UPDATE
LIMIT 1;

-- UpdateProduct :one
UPDATE products
    set product_name = $1,
    description = $2,
    price = $3,
    quantity = $4,
    discount = $5,
    rating = $6,
    size_options = $7,
    color_options = $8,
    category = $9,
    image_url = $10
WHERE id = $11
RETURNING *;

-- DeleteProduct :exec
DELETE FROM products
WHERE id = $1;