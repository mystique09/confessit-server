-- name: GetUserById :one
SELECT *
FROM "user"
WHERE id = $1
LIMIT 1;

-- name: GetUserByUsername :one
SELECT *
FROM "user"
WHERE username = $1
LIMIT 1;

-- name: ListUsers :many
SELECT *
FROM "user"
LIMIT 20
OFFSET $1;

-- name: CreateUser :one
INSERT INTO "user" (
    id, username, password
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: UpdateUsername :one
UPDATE "user"
SET username = $1
WHERE id = $2
RETURNING id;

-- name: DeleteOneUser :one
DELETE FROM "user"
WHERE id = $1
RETURNING id;