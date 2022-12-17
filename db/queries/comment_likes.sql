-- name: ListAllCommentLikes :many
SELECT * FROM "comment_likes" WHERE comment_id = $1;

-- name: CreateCommentLike :one
INSERT INTO "comment_likes" (id, comment_id, user_identity_id, type) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: DeleteCommentLike :one
DELETE FROM "comment_likes" WHERE comment_id = $1 AND user_identity_id = $2 RETURNING *;