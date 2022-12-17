-- name: ListAllPostLikes :many
SELECT * FROM "likes" WHERE post_id = $1;

-- name: CreatePostLike :one
INSERT INTO "likes" (id, post_id, user_identity_id) VALUES ($1, $2, $3) RETURNING *;

-- name: DeletePostLike :one
DELETE FROM "likes" WHERE post_id = $1 AND user_identity_id = $2 RETURNING *;