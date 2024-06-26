// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: reviews.sql

package db

import (
	"context"
)

const calculateProductRating = `-- name: CalculateProductRating :one
SELECT COALESCE(AVG(rating), 0) AS average_rating
FROM reviews
WHERE product_id = $1
`

func (q *Queries) CalculateProductRating(ctx context.Context, productID int64) (interface{}, error) {
	row := q.db.QueryRow(ctx, calculateProductRating, productID)
	var average_rating interface{}
	err := row.Scan(&average_rating)
	return average_rating, err
}

const createReveiw = `-- name: CreateReveiw :one
INSERT INTO reviews (
    user_id, product_id, rating, review
) VALUES (
    $1, $2, $3, $4
)
RETURNING id, user_id, product_id, rating, review
`

type CreateReveiwParams struct {
	UserID    int64  `json:"user_id"`
	ProductID int64  `json:"product_id"`
	Rating    int32  `json:"rating"`
	Review    string `json:"review"`
}

func (q *Queries) CreateReveiw(ctx context.Context, arg CreateReveiwParams) (Review, error) {
	row := q.db.QueryRow(ctx, createReveiw,
		arg.UserID,
		arg.ProductID,
		arg.Rating,
		arg.Review,
	)
	var i Review
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ProductID,
		&i.Rating,
		&i.Review,
	)
	return i, err
}

const deleteReview = `-- name: DeleteReview :exec
DELETE FROM reviews
WHERE id = $1
`

func (q *Queries) DeleteReview(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteReview, id)
	return err
}

const editReview = `-- name: EditReview :one
UPDATE reviews
    set rating = $1,
    review = $2
WHERE id = $3
RETURNING id, user_id, product_id, rating, review
`

type EditReviewParams struct {
	Rating int32  `json:"rating"`
	Review string `json:"review"`
	ID     int64  `json:"id"`
}

func (q *Queries) EditReview(ctx context.Context, arg EditReviewParams) (Review, error) {
	row := q.db.QueryRow(ctx, editReview, arg.Rating, arg.Review, arg.ID)
	var i Review
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ProductID,
		&i.Rating,
		&i.Review,
	)
	return i, err
}

const getProductReviews = `-- name: GetProductReviews :many
SELECT id, user_id, product_id, rating, review FROM reviews
WHERE product_id = $1
`

func (q *Queries) GetProductReviews(ctx context.Context, productID int64) ([]Review, error) {
	rows, err := q.db.Query(ctx, getProductReviews, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Review
	for rows.Next() {
		var i Review
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.ProductID,
			&i.Rating,
			&i.Review,
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

const getReview = `-- name: GetReview :one
SELECT id, user_id, product_id, rating, review FROM reviews
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetReview(ctx context.Context, id int64) (Review, error) {
	row := q.db.QueryRow(ctx, getReview, id)
	var i Review
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ProductID,
		&i.Rating,
		&i.Review,
	)
	return i, err
}

const getUsersReviews = `-- name: GetUsersReviews :many
SELECT id, user_id, product_id, rating, review FROM reviews
WHERE user_id = $1
`

func (q *Queries) GetUsersReviews(ctx context.Context, userID int64) ([]Review, error) {
	rows, err := q.db.Query(ctx, getUsersReviews, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Review
	for rows.Next() {
		var i Review
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.ProductID,
			&i.Rating,
			&i.Review,
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
