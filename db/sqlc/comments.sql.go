// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: comments.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createComment = `-- name: CreateComment :one
INSERT INTO "comments"(
	id,
 	content,
	user_identity_id, 
	post_id, 
	parent_id, 
	created_at, 
	updated_at
) VALUES (
	$1, $2, $3, $4, $5, $6, $7
) RETURNING id, content, user_identity_id, post_id, parent_id, created_at, updated_at
`

type CreateCommentParams struct {
	ID             uuid.UUID `json:"id"`
	Content        string    `json:"content"`
	UserIdentityID uuid.UUID `json:"user_identity_id"`
	PostID         uuid.UUID `json:"post_id"`
	ParentID       uuid.UUID `json:"parent_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (q *Queries) CreateComment(ctx context.Context, arg CreateCommentParams) (Comment, error) {
	row := q.queryRow(ctx, q.createCommentStmt, createComment,
		arg.ID,
		arg.Content,
		arg.UserIdentityID,
		arg.PostID,
		arg.ParentID,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Comment
	err := row.Scan(
		&i.ID,
		&i.Content,
		&i.UserIdentityID,
		&i.PostID,
		&i.ParentID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteComment = `-- name: DeleteComment :one
DELETE FROM "comments" WHERE id = $1 RETURNING id
`

func (q *Queries) DeleteComment(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	row := q.queryRow(ctx, q.deleteCommentStmt, deleteComment, id)
	err := row.Scan(&id)
	return id, err
}

const getComment = `-- name: GetComment :one
SELECT id, content, user_identity_id, post_id, parent_id, created_at, updated_at FROM "comments" WHERE id = $1
`

func (q *Queries) GetComment(ctx context.Context, id uuid.UUID) (Comment, error) {
	row := q.queryRow(ctx, q.getCommentStmt, getComment, id)
	var i Comment
	err := row.Scan(
		&i.ID,
		&i.Content,
		&i.UserIdentityID,
		&i.PostID,
		&i.ParentID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listAllComments = `-- name: ListAllComments :many
SELECT id, content, user_identity_id, post_id, parent_id, created_at, updated_at FROM "comments" WHERE post_id = $1
`

func (q *Queries) ListAllComments(ctx context.Context, postID uuid.UUID) ([]Comment, error) {
	rows, err := q.query(ctx, q.listAllCommentsStmt, listAllComments, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Comment
	for rows.Next() {
		var i Comment
		if err := rows.Scan(
			&i.ID,
			&i.Content,
			&i.UserIdentityID,
			&i.PostID,
			&i.ParentID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateComment = `-- name: UpdateComment :one
UPDATE "comments" SET content = $1, updated_at = $2 WHERE id = $3 RETURNING id
`

type UpdateCommentParams struct {
	Content   string    `json:"content"`
	UpdatedAt time.Time `json:"updated_at"`
	ID        uuid.UUID `json:"id"`
}

func (q *Queries) UpdateComment(ctx context.Context, arg UpdateCommentParams) (uuid.UUID, error) {
	row := q.queryRow(ctx, q.updateCommentStmt, updateComment, arg.Content, arg.UpdatedAt, arg.ID)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}
