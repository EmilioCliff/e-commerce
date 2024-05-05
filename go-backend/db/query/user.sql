-- name: CreateUser :one
INSERT INTO users (
    username, email, password, subscription, user_cart, role
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY username;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

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
    user_cart = $5,
    updated_at = $6
WHERE id = $7
RETURNING *;

-- name: UpdateUserCredentials :one
UPDATE users
    set username = $1,
    password = $2,
    role = $3,
    updated_at = $4
WHERE id = $4
RETURNING *;

-- name: UpdateUserSubscription :one
UPDATE users
    set subscription = $1,
    updated_at = $3
WHERE id = $2
RETURNING *;

-- name: UpdateUserCart :one
UPDATE users
    set user_cart = $1
WHERE id = $2
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
