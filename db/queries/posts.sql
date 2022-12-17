-- name: ListAllPosts :many
SELECT * FROM posts ORDER BY created_at DESC LIMIT 20 OFFSET $1;

-- name: GetPostById :one
SELECT * FROM posts WHERE id = $1 LIMIT 1;

-- name: CreatePost :one
INSERT INTO posts (id, content, user_identity_id) VALUES ($1, $2, $3) RETURNING *;

-- name: UpdatePost :one
UPDATE posts SET content = $1 WHERE id = $2 RETURNING id;

-- name: DeletePost :one
DELETE FROM posts WHERE id = $1 RETURNING id;