package handler

import (
	"cnfs/db/mock"
	db "cnfs/db/sqlc"
	"database/sql"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
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
	comment := RandomComment(t, post.ID, uuid.NullUUID{})

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
