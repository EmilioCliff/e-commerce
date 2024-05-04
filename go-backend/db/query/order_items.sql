-- name: CreateOrderItem :one
INSERT INTO order_items (
    order_id, product_id, color, size, quantity
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: ListOrderItems :many
SELECT * FROM order_items
WHERE id = $1;

-- name: GetOrderOrderItems :many
SELECT * FROM order_items
WHERE order_id = $1;

-- name: UpdateOrderItems :one
UPDATE order_items
    set color = $1,
    size = $2,
    quantity = $3
WHERE id = $4
RETURNING *;

-- name: DeleteOrderItem :exec
DELETE FROM order_items
WHERE id = $1;

-- name: GetOrderItem :one
SELECT * FROM order_items
WHERE id = $1;