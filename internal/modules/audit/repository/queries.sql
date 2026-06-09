-- name: GetSessionByRefreshToken :one
SELECT id, user_id, refresh_token, user_agent, ip_address, is_revoked, expires_at, created_at
FROM sessions
WHERE refresh_token = $1 LIMIT 1;

-- name: RevokeSession :exec
UPDATE sessions
SET is_revoked = TRUE
WHERE refresh_token = $1;