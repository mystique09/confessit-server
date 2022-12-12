-- name: GetUser :one
SELECT *
FROM "user"
WHERE id = $1
LIMIT 1;

-- name: ListUsers :many
SELECT *
FROM "user"
LIMIT 20
OFFSET $1;