-- name: CreateProduct :one
INSERT INTO products (
    product_name, description, price, quantity, discount, size_options, color_options, category, brand, image_url
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
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

-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1
LIMIT 1;

-- name: GetProductForUpdate :one
SELECT  * FROM products
WHERE id = $1
FOR NO KEY UPDATE
LIMIT 1;

-- name: UpdateProductRating :one
UPDATE products
    set rating = $2
WHERE id = $1
RETURNING *;

-- name: AddProductQuantity :one
UPDATE products
    set quantity = $1
WHERE id = $2
RETURNING *;

-- name: UpdateProduct :one
UPDATE products
    set product_name = $1,
    description = $2,
    price = $3,
    discount = $4,
    size_options = $5,
    color_options = $6,
    category = $7,
    brand = $8,
    image_url = $9,
    updated_at = $10
WHERE id = $11
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;