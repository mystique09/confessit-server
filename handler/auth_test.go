package handler

import (
	"cnfs/common"
	"cnfs/db/mock"
	db "cnfs/db/sqlc"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	password, user := RandomUser(t)

	testCases := []struct {
		name          string
		payload       string
		buildStubs    func(store *mock.MockStore)
		checkResponse func(rec *httptest.ResponseRecorder)
	}{
		{
			name:    "OK",
			payload: fmt.Sprintf(`{"username": %q, "password": %q}`, user.Username, password),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).Times(1).Return(user, nil)
				store.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Times(1)
				store.EXPECT().GetUserIdentityByUserId(gomock.Any(), gomock.Any()).Times(1)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 200, rec.Code)

				var resp response
				body, err := io.ReadAll(rec.Body)
				require.NoError(t, err)

				require.NoError(t, json.Unmarshal(body, &resp))
				require.NotNil(t, resp.Data)
			},
		},
		{
			name:    "Missing username field",
			payload: fmt.Sprintf(`{"password": %q}`, password),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 400, rec.Code)
				resp := new(response)

				body, err := io.ReadAll(rec.Body)
				require.NoError(t, err)
				require.NoError(t, json.Unmarshal(body, &resp))
				require.NotNil(t, resp.Err)
				require.Empty(t, resp.Data)
			},
		},
		{
			name:    "Missing password field",
			payload: fmt.Sprintf(`{"username": %q}`, user.Username),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 400, rec.Code)
				resp := new(response)

				body, err := io.ReadAll(rec.Body)
				require.NoError(t, err)
				require.NoError(t, json.Unmarshal(body, &resp))
				require.NotNil(t, resp.Err)
				require.Empty(t, resp.Data)
			},
		},
		{
			name:    "Missing payload",
			payload: ``,
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 400, rec.Code)
				resp := new(response)

				body, err := io.ReadAll(rec.Body)
				require.NoError(t, err)
				require.NoError(t, json.Unmarshal(body, &resp))
				require.NotNil(t, resp.Err)
				require.Empty(t, resp.Data)
			},
		},
		{
			name:    "NOT FOUND",
			payload: fmt.Sprintf(`{"username": %q, "password": %q}`, user.Username, password),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).Times(1).Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 404, rec.Code)
				var resp response

				body, err := io.ReadAll(rec.Body)
				require.NoError(t, err)

				require.NoError(t, json.Unmarshal(body, &resp))

				require.Equal(t, "user not found", resp.Err)
			},
		},
		{
			name:    "MISMATCH PASSWORD",
			payload: fmt.Sprintf(`{"username": %q, "password": %q}`, user.Username, common.RandomString(12)),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).Times(1).Return(user, nil)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 403, rec.Code)

				var resp response

				body, err := io.ReadAll(rec.Body)
				require.NoError(t, err)

				require.NoError(t, json.Unmarshal(body, &resp))

				require.Equal(t, "password mismatch", resp.Err)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock.NewMockStore(ctrl)
			tc.buildStubs(store)

			server, err := NewServer(store, cfg)
			require.NoError(t, err)

			rec := httptest.NewRecorder()
			payload := strings.NewReader(tc.payload)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth", payload)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			server.router.ServeHTTP(rec, req)
			tc.checkResponse(rec)
		})
	}
}

func TestClearSession(t *testing.T) {
	sessionId := uuid.New()

	testCases := []testCase{
		{
			name:    "OK",
			payload: fmt.Sprintf(`{"session_id": %q}`, sessionId),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().DeleteSession(gomock.Any(), gomock.Eq(sessionId)).Times(1).Return(sessionId, nil)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 200, rec.Code)
			},
		},
		{
			name:    "BAD REQUEST/Missing session ID",
			payload: "{}",
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().DeleteSession(gomock.Any(), gomock.Eq(sessionId)).Times(0)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 400, rec.Code)
			},
		},
		{
			name:    "SESSION NOT FOUND",
			payload: fmt.Sprintf(`{"session_id": %q}`, sessionId),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().DeleteSession(gomock.Any(), gomock.Eq(sessionId)).Times(1).Return(uuid.Nil, sql.ErrNoRows)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 404, rec.Code)
			},
		},
		{
			name:    "INTERNAL ERROR",
			payload: fmt.Sprintf(`{"session_id": %q}`, sessionId),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().DeleteSession(gomock.Any(), gomock.Eq(sessionId)).Times(1).Return(uuid.Nil, sql.ErrConnDone)
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

			req := httptest.NewRequest(http.MethodDelete, "/api/v1/auth/clear", strings.NewReader(tc.payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			server.router.ServeHTTP(rec, req)
			tc.checkResponse(rec)
		})
	}
}
