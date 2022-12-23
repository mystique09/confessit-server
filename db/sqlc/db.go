// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package db

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.blockSessionStmt, err = db.PrepareContext(ctx, blockSession); err != nil {
		return nil, fmt.Errorf("error preparing query BlockSession: %w", err)
	}
	if q.createCommentStmt, err = db.PrepareContext(ctx, createComment); err != nil {
		return nil, fmt.Errorf("error preparing query CreateComment: %w", err)
	}
	if q.createMessageStmt, err = db.PrepareContext(ctx, createMessage); err != nil {
		return nil, fmt.Errorf("error preparing query CreateMessage: %w", err)
	}
	if q.createPostStmt, err = db.PrepareContext(ctx, createPost); err != nil {
		return nil, fmt.Errorf("error preparing query CreatePost: %w", err)
	}
	if q.createSessionStmt, err = db.PrepareContext(ctx, createSession); err != nil {
		return nil, fmt.Errorf("error preparing query CreateSession: %w", err)
	}
	if q.createUserStmt, err = db.PrepareContext(ctx, createUser); err != nil {
		return nil, fmt.Errorf("error preparing query CreateUser: %w", err)
	}
	if q.createUserIdentityStmt, err = db.PrepareContext(ctx, createUserIdentity); err != nil {
		return nil, fmt.Errorf("error preparing query CreateUserIdentity: %w", err)
	}
	if q.deleteCommentStmt, err = db.PrepareContext(ctx, deleteComment); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteComment: %w", err)
	}
	if q.deleteOneMessageStmt, err = db.PrepareContext(ctx, deleteOneMessage); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteOneMessage: %w", err)
	}
	if q.deleteOneUserStmt, err = db.PrepareContext(ctx, deleteOneUser); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteOneUser: %w", err)
	}
	if q.deletePostStmt, err = db.PrepareContext(ctx, deletePost); err != nil {
		return nil, fmt.Errorf("error preparing query DeletePost: %w", err)
	}
	if q.deleteSessionStmt, err = db.PrepareContext(ctx, deleteSession); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteSession: %w", err)
	}
	if q.deleteSessionByUserIdStmt, err = db.PrepareContext(ctx, deleteSessionByUserId); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteSessionByUserId: %w", err)
	}
	if q.getCommentStmt, err = db.PrepareContext(ctx, getComment); err != nil {
		return nil, fmt.Errorf("error preparing query GetComment: %w", err)
	}
	if q.getMessageByIdStmt, err = db.PrepareContext(ctx, getMessageById); err != nil {
		return nil, fmt.Errorf("error preparing query GetMessageById: %w", err)
	}
	if q.getPostByIdStmt, err = db.PrepareContext(ctx, getPostById); err != nil {
		return nil, fmt.Errorf("error preparing query GetPostById: %w", err)
	}
	if q.getSessionByIdStmt, err = db.PrepareContext(ctx, getSessionById); err != nil {
		return nil, fmt.Errorf("error preparing query GetSessionById: %w", err)
	}
	if q.getUserByIdStmt, err = db.PrepareContext(ctx, getUserById); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserById: %w", err)
	}
	if q.getUserByUsernameStmt, err = db.PrepareContext(ctx, getUserByUsername); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserByUsername: %w", err)
	}
	if q.getUserIdentityByIdStmt, err = db.PrepareContext(ctx, getUserIdentityById); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserIdentityById: %w", err)
	}
	if q.getUserIdentityByUserIdStmt, err = db.PrepareContext(ctx, getUserIdentityByUserId); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserIdentityByUserId: %w", err)
	}
	if q.listAllCommentsStmt, err = db.PrepareContext(ctx, listAllComments); err != nil {
		return nil, fmt.Errorf("error preparing query ListAllComments: %w", err)
	}
	if q.listAllPostsStmt, err = db.PrepareContext(ctx, listAllPosts); err != nil {
		return nil, fmt.Errorf("error preparing query ListAllPosts: %w", err)
	}
	if q.listMessageStmt, err = db.PrepareContext(ctx, listMessage); err != nil {
		return nil, fmt.Errorf("error preparing query ListMessage: %w", err)
	}
	if q.listUsersStmt, err = db.PrepareContext(ctx, listUsers); err != nil {
		return nil, fmt.Errorf("error preparing query ListUsers: %w", err)
	}
	if q.updateCommentStmt, err = db.PrepareContext(ctx, updateComment); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateComment: %w", err)
	}
	if q.updateMessageStatusStmt, err = db.PrepareContext(ctx, updateMessageStatus); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateMessageStatus: %w", err)
	}
	if q.updatePostStmt, err = db.PrepareContext(ctx, updatePost); err != nil {
		return nil, fmt.Errorf("error preparing query UpdatePost: %w", err)
	}
	if q.updateUserPasswordStmt, err = db.PrepareContext(ctx, updateUserPassword); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateUserPassword: %w", err)
	}
	if q.updateUsernameStmt, err = db.PrepareContext(ctx, updateUsername); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateUsername: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.blockSessionStmt != nil {
		if cerr := q.blockSessionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing blockSessionStmt: %w", cerr)
		}
	}
	if q.createCommentStmt != nil {
		if cerr := q.createCommentStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createCommentStmt: %w", cerr)
		}
	}
	if q.createMessageStmt != nil {
		if cerr := q.createMessageStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createMessageStmt: %w", cerr)
		}
	}
	if q.createPostStmt != nil {
		if cerr := q.createPostStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createPostStmt: %w", cerr)
		}
	}
	if q.createSessionStmt != nil {
		if cerr := q.createSessionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createSessionStmt: %w", cerr)
		}
	}
	if q.createUserStmt != nil {
		if cerr := q.createUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createUserStmt: %w", cerr)
		}
	}
	if q.createUserIdentityStmt != nil {
		if cerr := q.createUserIdentityStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createUserIdentityStmt: %w", cerr)
		}
	}
	if q.deleteCommentStmt != nil {
		if cerr := q.deleteCommentStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteCommentStmt: %w", cerr)
		}
	}
	if q.deleteOneMessageStmt != nil {
		if cerr := q.deleteOneMessageStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteOneMessageStmt: %w", cerr)
		}
	}
	if q.deleteOneUserStmt != nil {
		if cerr := q.deleteOneUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteOneUserStmt: %w", cerr)
		}
	}
	if q.deletePostStmt != nil {
		if cerr := q.deletePostStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deletePostStmt: %w", cerr)
		}
	}
	if q.deleteSessionStmt != nil {
		if cerr := q.deleteSessionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteSessionStmt: %w", cerr)
		}
	}
	if q.deleteSessionByUserIdStmt != nil {
		if cerr := q.deleteSessionByUserIdStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteSessionByUserIdStmt: %w", cerr)
		}
	}
	if q.getCommentStmt != nil {
		if cerr := q.getCommentStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getCommentStmt: %w", cerr)
		}
	}
	if q.getMessageByIdStmt != nil {
		if cerr := q.getMessageByIdStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getMessageByIdStmt: %w", cerr)
		}
	}
	if q.getPostByIdStmt != nil {
		if cerr := q.getPostByIdStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPostByIdStmt: %w", cerr)
		}
	}
	if q.getSessionByIdStmt != nil {
		if cerr := q.getSessionByIdStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getSessionByIdStmt: %w", cerr)
		}
	}
	if q.getUserByIdStmt != nil {
		if cerr := q.getUserByIdStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserByIdStmt: %w", cerr)
		}
	}
	if q.getUserByUsernameStmt != nil {
		if cerr := q.getUserByUsernameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserByUsernameStmt: %w", cerr)
		}
	}
	if q.getUserIdentityByIdStmt != nil {
		if cerr := q.getUserIdentityByIdStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserIdentityByIdStmt: %w", cerr)
		}
	}
	if q.getUserIdentityByUserIdStmt != nil {
		if cerr := q.getUserIdentityByUserIdStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserIdentityByUserIdStmt: %w", cerr)
		}
	}
	if q.listAllCommentsStmt != nil {
		if cerr := q.listAllCommentsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listAllCommentsStmt: %w", cerr)
		}
	}
	if q.listAllPostsStmt != nil {
		if cerr := q.listAllPostsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listAllPostsStmt: %w", cerr)
		}
	}
	if q.listMessageStmt != nil {
		if cerr := q.listMessageStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listMessageStmt: %w", cerr)
		}
	}
	if q.listUsersStmt != nil {
		if cerr := q.listUsersStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listUsersStmt: %w", cerr)
		}
	}
	if q.updateCommentStmt != nil {
		if cerr := q.updateCommentStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateCommentStmt: %w", cerr)
		}
	}
	if q.updateMessageStatusStmt != nil {
		if cerr := q.updateMessageStatusStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateMessageStatusStmt: %w", cerr)
		}
	}
	if q.updatePostStmt != nil {
		if cerr := q.updatePostStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updatePostStmt: %w", cerr)
		}
	}
	if q.updateUserPasswordStmt != nil {
		if cerr := q.updateUserPasswordStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateUserPasswordStmt: %w", cerr)
		}
	}
	if q.updateUsernameStmt != nil {
		if cerr := q.updateUsernameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateUsernameStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                          DBTX
	tx                          *sql.Tx
	blockSessionStmt            *sql.Stmt
	createCommentStmt           *sql.Stmt
	createMessageStmt           *sql.Stmt
	createPostStmt              *sql.Stmt
	createSessionStmt           *sql.Stmt
	createUserStmt              *sql.Stmt
	createUserIdentityStmt      *sql.Stmt
	deleteCommentStmt           *sql.Stmt
	deleteOneMessageStmt        *sql.Stmt
	deleteOneUserStmt           *sql.Stmt
	deletePostStmt              *sql.Stmt
	deleteSessionStmt           *sql.Stmt
	deleteSessionByUserIdStmt   *sql.Stmt
	getCommentStmt              *sql.Stmt
	getMessageByIdStmt          *sql.Stmt
	getPostByIdStmt             *sql.Stmt
	getSessionByIdStmt          *sql.Stmt
	getUserByIdStmt             *sql.Stmt
	getUserByUsernameStmt       *sql.Stmt
	getUserIdentityByIdStmt     *sql.Stmt
	getUserIdentityByUserIdStmt *sql.Stmt
	listAllCommentsStmt         *sql.Stmt
	listAllPostsStmt            *sql.Stmt
	listMessageStmt             *sql.Stmt
	listUsersStmt               *sql.Stmt
	updateCommentStmt           *sql.Stmt
	updateMessageStatusStmt     *sql.Stmt
	updatePostStmt              *sql.Stmt
	updateUserPasswordStmt      *sql.Stmt
	updateUsernameStmt          *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                          tx,
		tx:                          tx,
		blockSessionStmt:            q.blockSessionStmt,
		createCommentStmt:           q.createCommentStmt,
		createMessageStmt:           q.createMessageStmt,
		createPostStmt:              q.createPostStmt,
		createSessionStmt:           q.createSessionStmt,
		createUserStmt:              q.createUserStmt,
		createUserIdentityStmt:      q.createUserIdentityStmt,
		deleteCommentStmt:           q.deleteCommentStmt,
		deleteOneMessageStmt:        q.deleteOneMessageStmt,
		deleteOneUserStmt:           q.deleteOneUserStmt,
		deletePostStmt:              q.deletePostStmt,
		deleteSessionStmt:           q.deleteSessionStmt,
		deleteSessionByUserIdStmt:   q.deleteSessionByUserIdStmt,
		getCommentStmt:              q.getCommentStmt,
		getMessageByIdStmt:          q.getMessageByIdStmt,
		getPostByIdStmt:             q.getPostByIdStmt,
		getSessionByIdStmt:          q.getSessionByIdStmt,
		getUserByIdStmt:             q.getUserByIdStmt,
		getUserByUsernameStmt:       q.getUserByUsernameStmt,
		getUserIdentityByIdStmt:     q.getUserIdentityByIdStmt,
		getUserIdentityByUserIdStmt: q.getUserIdentityByUserIdStmt,
		listAllCommentsStmt:         q.listAllCommentsStmt,
		listAllPostsStmt:            q.listAllPostsStmt,
		listMessageStmt:             q.listMessageStmt,
		listUsersStmt:               q.listUsersStmt,
		updateCommentStmt:           q.updateCommentStmt,
		updateMessageStatusStmt:     q.updateMessageStatusStmt,
		updatePostStmt:              q.updatePostStmt,
		updateUserPasswordStmt:      q.updateUserPasswordStmt,
		updateUsernameStmt:          q.updateUsernameStmt,
	}
}
