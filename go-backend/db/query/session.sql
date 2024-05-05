-- name: GetSession :one
SELECT * FROM sessions
WHERE id = $1 LIMIT 1;

-- name: CreateSession :one
INSERT INTO sessions (
    id, user_id, refresh_token, is_blocked, user_agent, user_ip, expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- nmae: UpdateRefreshToken :one
UPDATE sessions
    set refresh_token = $1,
    expires_at = $2,
    is_blocked = $3
WHERE user_id = $4
RETURNING *;

-- name: BlockSession :one
UPDATE sessions 
    set is_blocked = $1
WHERE id = $2
RETURNING *;