package handler

import (
	"cnfs/db/mock"
	db "cnfs/db/sqlc"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestListAllPublicPost(t *testing.T) {
	dummyPosts := make([]db.Post, 20)

	testCases := []testCase{
		{
			name:    "OK - with page 0",
			payload: "0",
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().ListAllPosts(gomock.Any(), gomock.Eq(int32(0))).Times(1).Return(dummyPosts, nil)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 200, rec.Code)
				resp := new(response)

				body, err := io.ReadAll(rec.Body)
				require.NoError(t, err)
				require.NoError(t, json.Unmarshal(body, &resp))
				require.NotNil(t, resp.Data)
			},
		},
		{
			name:    "OK - with page > 0",
			payload: "1",
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().ListAllPosts(gomock.Any(), gomock.Eq(int32(10))).Times(1).Return(dummyPosts, nil)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 200, rec.Code)
				resp := new(response)

				body, err := io.ReadAll(rec.Body)
				require.NoError(t, err)
				require.NoError(t, json.Unmarshal(body, &resp))
				require.NotNil(t, resp.Data)
			},
		},
		{
			name:    "NEGATIVE PAGE/INVALID PAGE",
			payload: "abc",
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().ListAllPosts(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 400, rec.Code)
				resp := new(response)

				body, err := io.ReadAll(rec.Body)
				require.NoError(t, err)
				require.NoError(t, json.Unmarshal(body, &resp))
				require.Nil(t, resp.Data)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock.NewMockStore(ctrl)
			tc.buildStubs(store)

			server, err := NewServer(store, *cfg)
			require.NoError(t, err)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1/posts?page="+tc.payload, nil)
			req.Header.Set("Content-Type", "application/json")

			server.router.ServeHTTP(rec, req)
			tc.checkResponse(rec)
		})
	}
}

type eqCreatePostParamsMatcher struct {
	arg            db.CreatePostParams
	Id             uuid.UUID
	UserIdentityID uuid.UUID
}

func (m *eqCreatePostParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreatePostParams)
	if !ok {
		return false
	}

	m.arg = arg
	m.Id = arg.ID
	m.UserIdentityID = arg.UserIdentityID
	return reflect.DeepEqual(arg, m.arg)
}

func (m *eqCreatePostParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v", m.arg)
}

func EqCreatePostParams(arg db.CreatePostParams) gomock.Matcher {
	return &eqCreatePostParamsMatcher{arg: arg}
}

func TestCreateNewPost(t *testing.T) {
	_, user := RandomUser(t)
	identityId := uuid.New()
	post := RandomPost(t, identityId)

	testCases := []testCase{
		{
			name:    "OK",
			payload: fmt.Sprintf(`{"content": %q, "user_identity_id": %q}`, post.Content, post.UserIdentityID),
			buildStubs: func(store *mock.MockStore) {
				arg := db.CreatePostParams{
					ID:             post.ID,
					Content:        post.Content,
					UserIdentityID: post.UserIdentityID,
				}
				store.EXPECT().CreatePost(gomock.Any(), EqCreatePostParams(arg)).Times(1).Return(post, nil)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 200, rec.Code)
				resp := new(response)

				body, err := io.ReadAll(rec.Body)
				require.NoError(t, err)
				require.NoError(t, json.Unmarshal(body, &resp))
				require.NotNil(t, resp.Data)
			},
		},
		{
			name:    "NEGATIVE - MISSING CONTENT",
			payload: fmt.Sprintf(`{"user_identity_id": %q}`, post.UserIdentityID),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().CreatePost(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 400, rec.Code)
				resp := new(response)

				body, err := io.ReadAll(rec.Body)
				require.NoError(t, err)
				require.NoError(t, json.Unmarshal(body, &resp))
				require.Nil(t, resp.Data)
			},
		},
		{
			name:    "NEGATIVE - USER IDENTITY ID NOT FOUND",
			payload: fmt.Sprintf(`{"content": %q, "user_identity_id": %q}`, post.Content, post.UserIdentityID),
			buildStubs: func(store *mock.MockStore) {
				arg := db.CreatePostParams{
					ID:             post.ID,
					Content:        post.Content,
					UserIdentityID: post.UserIdentityID,
				}
				store.EXPECT().CreatePost(gomock.Any(), EqCreatePostParams(arg)).Times(1).Return(db.Post{}, sql.ErrNoRows)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 404, rec.Code)
				resp := new(response)

				body, err := io.ReadAll(rec.Body)
				require.NoError(t, err)
				require.NoError(t, json.Unmarshal(body, &resp))
				require.Nil(t, resp.Data)
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

			server, err := NewServer(store, *cfg)
			require.NoError(t, err)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/v1/posts", strings.NewReader(tc.payload))
			req.Header.Set("Content-Type", "application/json")
			token, _, err := server.tokenMaker.CreateToken(user.ID, user.Username, server.tokenCfg.GetAccessTokenDuration())
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+token)

			server.router.ServeHTTP(rec, req)
			tc.checkResponse(rec)
		})
	}
}

func TestGetPostById(t *testing.T) {
	post := RandomPost(t, uuid.New())

	testCases := []testCase{
		{
			name: "OK",
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetPostById(gomock.Any(), post.ID).Times(1).Return(post, nil)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 200, rec.Code)
				resp := new(response)

				body, err := io.ReadAll(rec.Body)
				require.NoError(t, err)
				require.NoError(t, json.Unmarshal(body, &resp))
				require.NotNil(t, resp.Data)
			},
		},
		{
			name: "NEGATIVE - POST NOT FOUND",
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetPostById(gomock.Any(), post.ID).Times(1).Return(db.Post{}, sql.ErrNoRows)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 404, rec.Code)
				resp := new(response)

				body, err := io.ReadAll(rec.Body)
				require.NoError(t, err)
				require.NoError(t, json.Unmarshal(body, &resp))
				require.Nil(t, resp.Data)
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

			server, err := NewServer(store, *cfg)
			require.NoError(t, err)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1/posts/"+post.ID.String(), nil)
			req.Header.Set("Content-Type", "application/json")

			server.router.ServeHTTP(rec, req)
			tc.checkResponse(rec)
		})
	}
}

func TestUpdateOnePost(t *testing.T) {
	userIdentityId := uuid.New()
	post := RandomPost(t, userIdentityId)
	_, user := RandomUser(t)

	testCases := []struct {
		name          string
		payload       string
		buildStubs    func(store *mock.MockStore)
		checkResponse func(rec *httptest.ResponseRecorder)
	}{
		{
			name:    "OK",
			payload: fmt.Sprintf(`{"content": %q}`, post.Content),
			buildStubs: func(store *mock.MockStore) {
				arg := db.UpdatePostParams{
					ID:      post.ID,
					Content: post.Content,
				}

				store.EXPECT().GetUserIdentityByUserId(gomock.Any(), gomock.Eq(user.ID)).Return(db.UserIdentity{
					ID:           userIdentityId,
					UserID:       user.ID,
					IdentityHash: uuid.New(),
				}, nil)

				store.EXPECT().GetPostById(gomock.Any(), gomock.Eq(post.ID)).Return(post, nil)
				store.EXPECT().UpdatePost(gomock.Any(), gomock.Eq(arg)).Return(post.ID, nil)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 200, rec.Code)
			},
		},
		{
			name:    "NEGATIVE - POST NOT FOUND",
			payload: fmt.Sprintf(`{"content": %q}`, post.Content),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetUserIdentityByUserId(gomock.Any(), gomock.Eq(user.ID)).Times(1).Return(db.UserIdentity{
					ID:           userIdentityId,
					UserID:       user.ID,
					IdentityHash: uuid.New(),
				}, nil)

				store.EXPECT().GetPostById(gomock.Any(), gomock.Eq(post.ID)).Times(1).Return(db.Post{}, sql.ErrNoRows)
				store.EXPECT().UpdatePost(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 404, rec.Code)
			},
		},
		{
			name:    "Unautorized",
			payload: fmt.Sprintf(`{"content": %q}`, post.Content),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetUserIdentityByUserId(gomock.Any(), gomock.Eq(user.ID)).Times(1).Return(db.UserIdentity{
					ID:           uuid.New(),
					UserID:       user.ID,
					IdentityHash: uuid.New(),
				}, nil)
				store.EXPECT().GetPostById(gomock.Any(), gomock.Eq(post.ID)).Times(1).Return(post, nil)
				store.EXPECT().UpdatePost(gomock.Any(), gomock.Any()).Times(0)
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

			server, err := NewServer(store, *cfg)
			require.NoError(t, err)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("PATCH", "/api/v1/posts/"+post.ID.String(), strings.NewReader(tc.payload))
			req.Header.Set("Content-Type", "application/json")
			token, _, err := server.tokenMaker.CreateToken(user.ID, user.Username, server.tokenCfg.GetAccessTokenDuration())
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+token)

			server.router.ServeHTTP(rec, req)
			tc.checkResponse(rec)
		})
	}
}

func TestDeletePost(t *testing.T) {
	userIdentityId := uuid.New()
	post := RandomPost(t, userIdentityId)
	_, user := RandomUser(t)

	testCases := []testCase{
		{
			name:    "OK",
			payload: "/api/v1/posts/" + post.ID.String(),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetUserIdentityByUserId(gomock.Any(), gomock.Eq(user.ID)).Times(1).Return(db.UserIdentity{
					ID:           userIdentityId,
					UserID:       user.ID,
					IdentityHash: uuid.New(),
				}, nil)

				store.EXPECT().GetPostById(gomock.Any(), gomock.Eq(post.ID)).Times(1).Return(post, nil)
				store.EXPECT().DeletePost(gomock.Any(), gomock.Eq(post.ID)).Times(1).Return(post.ID, nil)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 200, rec.Code)
			},
		},
		{
			name:    "NEGATIVE - POST NOT FOUND",
			payload: "/api/v1/posts/" + post.ID.String(),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetUserIdentityByUserId(gomock.Any(), gomock.Eq(user.ID)).Times(1).Return(db.UserIdentity{
					ID:           userIdentityId,
					UserID:       user.ID,
					IdentityHash: uuid.New(),
				}, nil)

				store.EXPECT().GetPostById(gomock.Any(), gomock.Eq(post.ID)).Times(1).Return(db.Post{}, sql.ErrNoRows)
				store.EXPECT().DeletePost(gomock.Any(), gomock.Eq(post.ID)).Times(0)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 404, rec.Code)
			},
		},
		{
			name:    "Unautorized",
			payload: "/api/v1/posts/" + post.ID.String(),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetUserIdentityByUserId(gomock.Any(), gomock.Eq(user.ID)).Times(1).Return(db.UserIdentity{
					ID:           uuid.New(),
					UserID:       user.ID,
					IdentityHash: uuid.New(),
				}, nil)
				store.EXPECT().GetPostById(gomock.Any(), gomock.Eq(post.ID)).Times(1).Return(post, nil)
				store.EXPECT().DeletePost(gomock.Any(), gomock.Eq(post.ID)).Times(0)
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

			server, err := NewServer(store, *cfg)
			require.NoError(t, err)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", tc.payload, nil)
			req.Header.Set("Content-Type", "application/json")
			token, _, err := server.tokenMaker.CreateToken(user.ID, user.Username, server.tokenCfg.GetAccessTokenDuration())
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+token)

			server.router.ServeHTTP(rec, req)
			tc.checkResponse(rec)
		})
	}
}
