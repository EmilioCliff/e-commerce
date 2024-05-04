-- name: CreateReveiw :one
INSERT INTO reviews (
    user_id, product_id, rating, review
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetProductReviews :many
SELECT * FROM reviews
WHERE product_id = $1;

-- name: DeleteReview :exec
DELETE FROM reviews
WHERE id = $1;

-- name: EditReview :one
UPDATE reviews
    set rating = $1,
    review = $2
WHERE id = $3
RETURNING *;

-- name: GetReview :one
SELECT * FROM reviews
WHERE id = $1 LIMIT 1;

-- name: GetUsersReviews :many
SELECT * FROM reviews
WHERE user_id = $1;