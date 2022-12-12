-- name: GetMessageById :one
SELECT *
FROM "message"
WHERE id = $1 AND receiver_id = $2
LIMIT 1;

-- name: ListMessage :many
SELECT *
FROM "message"
WHERE id = $1 AND receiver_id = $2 
LIMIT 20
OFFSET $3;

-- name: CreateMessage :one
INSERT INTO "message" (
    id, receiver_id, content
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: UpdateMessageStatus :one
UPDATE "message"
SET seen = TRUE
WHERE id = $1 AND receiver_id = $2
RETURNING id;

-- name: DeleteOneMessage :one
DELETE FROM "message"
WHERE id = $1 AND receiver_id = $2
RETURNING id;