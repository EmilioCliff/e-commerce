-- name: CreateOrder :one
INSERT INTO orders (
    user_id, amount, status, shipping_address, created_at
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetOrder :one
SELECT * FROM orders
WHERE id = $1 LIMIT 1;

-- name: GetOrderForUpdate :one
SELECT * FROM orders
WHERE id = $1
FOR NO KEY UPDATE;

-- name: GetUserOrders :many
SELECT * FROM orders
WHERE user_id = $1;

-- name: UpdateOrder :one
UPDATE orders
    set status = $1
WHERE id = $2
RETURNING *;

-- name: ListOrders :many
SELECT * FROM orders
ORDER BY created_at DESC;

-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1;