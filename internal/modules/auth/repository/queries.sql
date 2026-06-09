-- name: CreateSession :one
INSERT INTO sessions (user_id, refresh_token, expires_at)
VALUES ($1, $2, $3)
RETURNING id, user_id, refresh_token, is_revoked, expires_at, created_at;

-- name: GetSessionByToken :one
SELECT id, user_id, refresh_token, is_revoked, expires_at, created_at
FROM sessions
WHERE refresh_token = $1 LIMIT 1;

-- name: RevokeSession :exec
UPDATE sessions
SET is_revoked = TRUE
WHERE id = $1;

-- name: RevokeAllUserSessions :exec
UPDATE sessions
SET is_revoked = TRUE
WHERE user_id = $1;

-- name: CreatePasswordReset :one
INSERT INTO password_resets (user_id, token, expires_at)
VALUES ($1, $2, $3)
RETURNING id, user_id, token, used, expires_at, created_at;

-- name: GetPasswordResetByToken :one
SELECT id, user_id, token, used, expires_at, created_at
FROM password_resets
WHERE token = $1 LIMIT 1;

-- name: MarkPasswordResetUsed :exec
UPDATE password_resets
SET used = TRUE
WHERE id = $1;
