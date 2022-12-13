-- name: GetSessionById :one
SELECT *
FROM "session"
WHERE id = $1
LIMIT 1;

-- name: CreateSession :one
INSERT INTO "session"(
    id,
    user_id,
    username,
    refresh_token,
    user_agent,
    client_ip,
    is_blocked,
    expires_at 
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: BlockSession :one
UPDATE "session"
SET is_blocked = true
WHERE id = $1
RETURNING id;

-- name: DeleteSession :one
DELETE FROM "session"
WHERE id = $1
RETURNING id;

-- name: DeleteSessionByUserId :one
DELETE FROM "session"
WHERE user_id = $1
RETURNING id;