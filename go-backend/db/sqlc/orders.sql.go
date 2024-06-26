// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: orders.sql

package db

import (
	"context"
)

const createOrder = `-- name: CreateOrder :one
INSERT INTO orders (
    user_id, amount, shipping_address
) VALUES (
    $1, $2, $3
)
RETURNING id, user_id, amount, status, shipping_address, created_at
`

type CreateOrderParams struct {
	UserID          int64   `json:"user_id"`
	Amount          float64 `json:"amount"`
	ShippingAddress string  `json:"shipping_address"`
}

func (q *Queries) CreateOrder(ctx context.Context, arg CreateOrderParams) (Order, error) {
	row := q.db.QueryRow(ctx, createOrder, arg.UserID, arg.Amount, arg.ShippingAddress)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Amount,
		&i.Status,
		&i.ShippingAddress,
		&i.CreatedAt,
	)
	return i, err
}

const deleteOrder = `-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1
`

func (q *Queries) DeleteOrder(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteOrder, id)
	return err
}

const getOrder = `-- name: GetOrder :one
SELECT id, user_id, amount, status, shipping_address, created_at FROM orders
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetOrder(ctx context.Context, id int64) (Order, error) {
	row := q.db.QueryRow(ctx, getOrder, id)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Amount,
		&i.Status,
		&i.ShippingAddress,
		&i.CreatedAt,
	)
	return i, err
}

const getOrderForUpdate = `-- name: GetOrderForUpdate :one
SELECT id, user_id, amount, status, shipping_address, created_at FROM orders
WHERE id = $1
FOR NO KEY UPDATE
`

func (q *Queries) GetOrderForUpdate(ctx context.Context, id int64) (Order, error) {
	row := q.db.QueryRow(ctx, getOrderForUpdate, id)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Amount,
		&i.Status,
		&i.ShippingAddress,
		&i.CreatedAt,
	)
	return i, err
}

const getUserOrders = `-- name: GetUserOrders :many
SELECT id, user_id, amount, status, shipping_address, created_at FROM orders
WHERE user_id = $1
`

func (q *Queries) GetUserOrders(ctx context.Context, userID int64) ([]Order, error) {
	rows, err := q.db.Query(ctx, getUserOrders, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Order
	for rows.Next() {
		var i Order
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Amount,
			&i.Status,
			&i.ShippingAddress,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listOrders = `-- name: ListOrders :many
SELECT id, user_id, amount, status, shipping_address, created_at FROM orders
ORDER BY created_at DESC
`

func (q *Queries) ListOrders(ctx context.Context) ([]Order, error) {
	rows, err := q.db.Query(ctx, listOrders)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Order
	for rows.Next() {
		var i Order
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Amount,
			&i.Status,
			&i.ShippingAddress,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateOrder = `-- name: UpdateOrder :one
UPDATE orders
    set status = $1
WHERE id = $2
RETURNING id, user_id, amount, status, shipping_address, created_at
`

type UpdateOrderParams struct {
	Status string `json:"status"`
	ID     int64  `json:"id"`
}

func (q *Queries) UpdateOrder(ctx context.Context, arg UpdateOrderParams) (Order, error) {
	row := q.db.QueryRow(ctx, updateOrder, arg.Status, arg.ID)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Amount,
		&i.Status,
		&i.ShippingAddress,
		&i.CreatedAt,
	)
	return i, err
}
