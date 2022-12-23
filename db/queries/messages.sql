-- name: GetMessageById :one
SELECT *
FROM "messages"
WHERE id = $1
LIMIT 1;

-- name: ListMessage :many
SELECT *
FROM "messages"
WHERE receiver_id = $1
LIMIT 20
OFFSET $2;

-- name: CreateMessage :one
INSERT INTO "messages" (
    id, receiver_id, content, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: UpdateMessageStatus :one
UPDATE "messages"
SET seen = TRUE, updated_at = $1
WHERE id = $2 AND receiver_id = $3
RETURNING id;

-- name: DeleteOneMessage :one
DELETE FROM "messages"
WHERE id = $1 AND receiver_id = $2
RETURNING id;