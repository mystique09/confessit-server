package handler

import (
	"cnfs/common"
	"cnfs/db/mock"
	db "cnfs/db/sqlc"
	"database/sql"
	"fmt"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestListAllComments(t *testing.T) {
	comments := make([]db.Comment, 20)
	post := RandomPost(t, uuid.New())

	testCases := []testCase{
		{
			name:    "OK",
			payload: "/api/v1/posts/" + post.ID.String() + "/comments",
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().ListAllComments(gomock.Any(), gomock.Eq(post.ID)).Return(comments, nil)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 200, rec.Code)
			},
		},
		{
			name:    "404 post not found",
			payload: "/api/v1/posts/" + post.ID.String() + "/comments",
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().ListAllComments(gomock.Any(), gomock.Eq(post.ID)).Return([]db.Comment{}, sql.ErrNoRows)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 404, rec.Code)
			},
		},
		{
			name:    "500 internal server error",
			payload: "/api/v1/posts/" + post.ID.String() + "/comments",
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().ListAllComments(gomock.Any(), gomock.Eq(post.ID)).Return([]db.Comment{}, sql.ErrConnDone)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 500, rec.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock.NewMockStore(ctrl)
			tc.buildStubs(store)

			server, err := NewServer(store, cfg)
			require.NoError(t, err)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", tc.payload, nil)

			server.router.ServeHTTP(rec, req)
			tc.checkResponse(rec)
		})
	}
}

func TestGetComment(t *testing.T) {
	post := RandomPost(t, uuid.New())
	comment := RandomComment(t, post.ID, uuid.Nil)

	testCases := []testCase{
		{
			name:    "OK",
			payload: "/api/v1/comments/" + comment.ID.String(),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetComment(gomock.Any(), gomock.Eq(comment.ID)).Return(comment, nil)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 200, rec.Code)
			},
		},
		{
			name:    "404 comment not found",
			payload: "/api/v1/comments/" + comment.ID.String(),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetComment(gomock.Any(), gomock.Eq(comment.ID)).Return(db.Comment{}, sql.ErrNoRows)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 404, rec.Code)
			},
		},
		{
			name:    "500 internal server error",
			payload: "/api/v1/comments/" + comment.ID.String(),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetComment(gomock.Any(), gomock.Eq(comment.ID)).Return(db.Comment{}, sql.ErrConnDone)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 500, rec.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock.NewMockStore(ctrl)
			tc.buildStubs(store)

			server, err := NewServer(store, cfg)
			require.NoError(t, err)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", tc.payload, nil)

			server.router.ServeHTTP(rec, req)
			tc.checkResponse(rec)
		})
	}
}

type eqCreateCommentParams struct {
	arg db.CreateCommentParams
	ID  uuid.UUID
}

func (e *eqCreateCommentParams) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateCommentParams)
	if !ok {
		return false
	}

	arg.ID = e.ID
	arg = e.arg
	return reflect.DeepEqual(e.arg, arg)
}

func (e *eqCreateCommentParams) String() string {
	return fmt.Sprintf("is equal to %v", e.arg)
}

func EqCreateCommentParams(arg *db.CreateCommentParams, id uuid.UUID) gomock.Matcher {
	return &eqCreateCommentParams{arg: *arg, ID: id}
}

func TestCreateComment(t *testing.T) {
	_, user := RandomUser(t)
	post := RandomPost(t, uuid.New())
	comment := RandomComment(t, post.ID, uuid.Nil)
	arg := db.CreateCommentParams{
		ID:             comment.ID,
		PostID:         post.ID,
		Content:        comment.Content,
		UserIdentityID: comment.UserIdentityID,
		ParentID:       comment.ParentID,
		CreatedAt:      comment.CreatedAt,
		UpdatedAt:      comment.UpdatedAt,
	}

	testCases := []testCase{
		{
			name:    "OK",
			payload: fmt.Sprintf(`{"user_identity_id": %q, "post_id": %q, "parent_id": %q, "content": %q}`, comment.UserIdentityID, post.ID, comment.ParentID, comment.Content),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().CreateComment(gomock.Any(), EqCreateCommentParams(&arg, arg.ID)).Times(1).Return(comment, nil)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 200, rec.Code)
			},
		},
		{
			name:    "404 post not found",
			payload: fmt.Sprintf(`{"user_identity_id": %q, "post_id": %q, "parent_id": %q, "content": %q}`, comment.UserIdentityID, post.ID, comment.ParentID, comment.Content),

			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().CreateComment(gomock.Any(), EqCreateCommentParams(&arg, arg.ID)).Times(1).Return(db.Comment{}, sql.ErrNoRows)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 404, rec.Code)
			},
		},
		{
			name:    "500 internal server error",
			payload: fmt.Sprintf(`{"user_identity_id": %q, "post_id": %q, "parent_id": %q, "content": %q}`, comment.UserIdentityID, post.ID, comment.ParentID, comment.Content),

			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().CreateComment(gomock.Any(), EqCreateCommentParams(&arg, arg.ID)).Times(1).Return(db.Comment{}, sql.ErrConnDone)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 500, rec.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock.NewMockStore(ctrl)
			tc.buildStubs(store)

			server, err := NewServer(store, cfg)
			require.NoError(t, err)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/v1/comments", strings.NewReader(tc.payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			tokenPayload, _, err := server.tokenMaker.CreateToken(user.ID, user.Username, cfg.AccessTokenDuration)
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+tokenPayload)

			server.router.ServeHTTP(rec, req)
			tc.checkResponse(rec)
		})
	}
}

type eqUpdateCommentParams struct {
	arg       db.UpdateCommentParams
	updatedAt time.Time
}

func (e *eqUpdateCommentParams) Matches(x interface{}) bool {
	arg, ok := x.(db.UpdateCommentParams)
	if !ok {
		return false
	}

	arg.UpdatedAt = e.updatedAt
	// arg = e.arg
	return reflect.DeepEqual(e.arg, arg)
}

func (e *eqUpdateCommentParams) String() string {
	return fmt.Sprintf("is equal to %v", e.arg)
}

func EqUpdateCommentParams(arg db.UpdateCommentParams, updatedAt time.Time) gomock.Matcher {
	return &eqUpdateCommentParams{arg: arg, updatedAt: updatedAt}
}

func TestUpdateComment(t *testing.T) {
	_, user := RandomUser(t)
	post := RandomPost(t, uuid.New())
	comment := RandomComment(t, post.ID, uuid.Nil)
	newContent := common.RandomString(36)

	testCases := []testCase{
		{
			name:    "OK",
			payload: fmt.Sprintf(`{"content": %q}`, newContent),
			buildStubs: func(store *mock.MockStore) {
				updatedAt := time.Now()
				arg := db.UpdateCommentParams{
					UpdatedAt: updatedAt,
					Content:   newContent,
					ID:        comment.ID,
				}

				store.EXPECT().GetUserIdentityByUserId(gomock.Any(), gomock.Eq(user.ID)).Times(1)
				store.EXPECT().GetComment(gomock.Any(), gomock.Eq(comment.ID)).Times(1)
				comment.UpdatedAt = updatedAt
				store.EXPECT().UpdateComment(gomock.Any(), EqUpdateCommentParams(arg, updatedAt)).Times(1).Return(comment.ID, nil)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 200, rec.Code)
			},
		},
		{
			name:    "400 bad request - Missing payload",
			payload: "",
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetUserIdentityByUserId(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 400, rec.Code)
			},
		},
		{
			name:    "Unauthorized - Cannot update a comment that does not belong to the user",
			payload: `{"content": "new content"}`,
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetUserIdentityByUserId(gomock.Any(), gomock.Eq(user.ID)).Times(1).Return(db.UserIdentity{
					ID:           uuid.New(),
					UserID:       uuid.New(),
					IdentityHash: uuid.New(),
				}, nil)
				store.EXPECT().GetComment(gomock.Any(), gomock.Eq(comment.ID)).Times(1).Return(comment, nil)
				store.EXPECT().UpdateComment(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 401, rec.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock.NewMockStore(ctrl)
			tc.buildStubs(store)

			server, err := NewServer(store, cfg)
			require.NoError(t, err)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/api/v1/comments/"+comment.ID.String(), strings.NewReader(tc.payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			tokenPayload, _, err := server.tokenMaker.CreateToken(user.ID, user.Username, cfg.AccessTokenDuration)
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+tokenPayload)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			server.router.ServeHTTP(rec, req)
			tc.checkResponse(rec)
		})
	}
}

func TestDeleteComment(t *testing.T) {
	_, user := RandomUser(t)
	post := RandomPost(t, uuid.New())
	comment := RandomComment(t, post.ID, uuid.Nil)

	testCases := []testCase{
		{
			name:    "OK",
			payload: comment.ID.String(),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetUserIdentityByUserId(gomock.Any(), gomock.Eq(user.ID)).Times(1)
				store.EXPECT().GetComment(gomock.Any(), gomock.Eq(comment.ID)).Times(1)
				store.EXPECT().DeleteComment(gomock.Any(), gomock.Eq(comment.ID)).Times(1).Return(comment.ID, nil)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 200, rec.Code)
			},
		},
		{
			name:    "Unauthorized - Cannot delete a comment that does not belong to the user",
			payload: comment.ID.String(),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetUserIdentityByUserId(gomock.Any(), gomock.Eq(user.ID)).Times(1).Return(db.UserIdentity{
					ID:           uuid.New(),
					UserID:       uuid.New(),
					IdentityHash: uuid.New(),
				}, nil)
				store.EXPECT().GetComment(gomock.Any(), gomock.Eq(comment.ID)).Times(1).Return(comment, nil)
				store.EXPECT().DeleteComment(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 401, rec.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock.NewMockStore(ctrl)
			tc.buildStubs(store)

			server, err := NewServer(store, cfg)
			require.NoError(t, err)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/api/v1/comments/"+tc.payload, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			tokenPayload, _, err := server.tokenMaker.CreateToken(user.ID, user.Username, cfg.AccessTokenDuration)
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+tokenPayload)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			server.router.ServeHTTP(rec, req)
			tc.checkResponse(rec)
		})
	}
}
