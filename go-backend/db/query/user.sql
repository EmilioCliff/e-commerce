-- name: CreateUser :one
INSERT INTO users (
    username, email, password, subscription, token, refresh_token, user_cart, role
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY username;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserForUpdate :one
SELECT * FROM users
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: GetSubscribedUsers :many
SELECT * FROM users
WHERE
    subscription = true
ORDER BY username;

-- name: UpdateUser :one
UPDATE users
    set username = $1,
    email = $2,
    password = $3,
    subscription = $4,
    token = $5,
    refresh_token = $6,
    user_cart = $7,
    updated_at = $8
WHERE id = $9
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
