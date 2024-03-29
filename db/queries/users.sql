-- name: GetUserById :one
SELECT *
FROM "users"
WHERE id = $1
LIMIT 1;

-- name: GetUserByUsername :one
SELECT *
FROM "users"
WHERE username = $1
LIMIT 1;

-- name: ListUsers :many
SELECT *
FROM "users"
LIMIT 20
OFFSET $1;

-- name: CreateUser :one
INSERT INTO "users" (
    id, username, password, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING id;

-- name: UpdateUsername :one
UPDATE "users"
SET username = $1, updated_at = $2
WHERE id = $3
RETURNING id;

-- name: UpdateUserPassword :one
UPDATE "users"
SET password = $1, updated_at = $2
WHERE id = $3
RETURNING id;

-- name: DeleteOneUser :one
DELETE FROM "users"
WHERE id = $1
RETURNING id;