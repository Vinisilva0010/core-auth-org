-- name: CreateSession :one
INSERT INTO sessions (
    user_id, refresh_token, user_agent, ip_address, expires_at
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING id, user_id, refresh_token, user_agent, ip_address, is_revoked, expires_at, created_at;

-- name: GetSession :one
SELECT id, user_id, refresh_token, user_agent, ip_address, is_revoked, expires_at, created_at
FROM sessions
WHERE id = $1 LIMIT 1;

-- name: GetSessionByRefreshToken :one
SELECT id, user_id, refresh_token, user_agent, ip_address, is_revoked, expires_at, created_at
FROM sessions
WHERE refresh_token = $1 LIMIT 1;

-- name: RevokeSession :exec
UPDATE sessions
SET is_revoked = TRUE
WHERE refresh_token = $1;
