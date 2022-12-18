-- name: ListAllComments :many
SELECT * FROM "comments" WHERE post_id = $1;

-- name: GetComment :one
SELECT * FROM "comments" WHERE id = $1;

-- name: CreateComment :one
INSERT INTO "comments" (id, content, user_identity_id, post_id, parent_id) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UpdateComment :one
UPDATE "comments" SET content = $1, updated_at = $2 WHERE id = $3 RETURNING *;

-- name: DeleteComment :one
DELETE FROM "comments" WHERE id = $1 RETURNING *;