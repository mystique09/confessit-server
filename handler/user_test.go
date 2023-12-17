package handler

import (
	"cnfs/common"
	"cnfs/db/mock"
	db "cnfs/db/sqlc"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e *eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := common.CheckPassword([]byte(arg.Password), []byte(e.password))
	if err != nil {
		return false
	}

	e.arg.ID = arg.ID
	e.arg.Password = arg.Password
	e.arg.CreatedAt = arg.CreatedAt
	e.arg.UpdatedAt = arg.UpdatedAt

	return reflect.DeepEqual(e.arg, arg)
}

func (e *eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and id %v", e.arg, e.password)
}

func EqCreateUserParams(arg *db.CreateUserParams, password string) gomock.Matcher {
	return &eqCreateUserParamsMatcher{*arg, password}
}

func TestCreateUser(t *testing.T) {
	password, user := RandomUser(t)

	testCases := []testCase{
		{
			name:    "OK",
			payload: fmt.Sprintf(`{"username": %q, "password": %q}`, user.Username, password),
			buildStubs: func(store *mock.MockStore) {
				arg := &db.CreateUserParams{
					ID:        user.ID,
					Username:  user.Username,
					Password:  user.Password,
					CreatedAt: user.CreatedAt,
					UpdatedAt: user.UpdatedAt,
				}

				store.
					EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(user.ID, nil)
				store.EXPECT().CreateUserIdentity(gomock.Any(), gomock.Any()).Times(1)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 200, rec.Code)
				resp := new(response)

				body, err := io.ReadAll(rec.Body)
				require.NoError(t, err)
				require.NoError(t, json.Unmarshal(body, &resp))
				require.NotNil(t, resp.Data)
				require.Empty(t, resp.Err)
			},
		},
		{
			name:    "Missing field",
			payload: `{"password":"testpassword"}`,
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 400, rec.Code)

				resp := new(response)

				body, err := io.ReadAll(rec.Body)
				require.NoError(t, err)

				require.NoError(t, json.Unmarshal(body, &resp))

				require.Contains(t, resp.Err, "Username", "required")
			},
		},
		{
			name:    "No payload",
			payload: ``,
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 400, rec.Code)

				resp := new(response)

				body, err := io.ReadAll(rec.Body)
				require.NoError(t, err)

				require.NoError(t, json.Unmarshal(body, &resp))

				require.Contains(t, resp.Err, "Password", "Username", "required")
			},
		},
		{
			name:    "Already exist",
			payload: fmt.Sprintf(`{"username": %q, "password": %q}`, user.Username, password),
			buildStubs: func(store *mock.MockStore) {
				arg := &db.CreateUserParams{
					ID:        user.ID,
					Username:  user.Username,
					Password:  user.Password,
					CreatedAt: user.CreatedAt,
					UpdatedAt: user.UpdatedAt,
				}

				store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).Times(1).Return(user.ID, errors.New("unique violation"))
				store.EXPECT().CreateUserIdentity(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 400, rec.Code)

				resp := new(response)

				body, err := io.ReadAll(rec.Body)
				require.NoError(t, err)

				require.NoError(t, json.Unmarshal(body, &resp))
				require.Equal(t, "user already exist", resp.Err)
			},
		},
		{
			name:    "INTERNAL ERROR",
			payload: fmt.Sprintf(`{"username": %q, "password": %q}`, user.Username, password),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(uuid.Nil, sql.ErrConnDone)
				store.EXPECT().CreateUserIdentity(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 500, rec.Code)

				resp := new(response)

				body, err := io.ReadAll(rec.Body)
				require.NoError(t, err)

				require.NoError(t, json.Unmarshal(body, &resp))
				require.NotNil(t, resp.Err)
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
			payload := strings.NewReader(tc.payload)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/users", payload)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			server.router.ServeHTTP(rec, req)
			tc.checkResponse(rec)
		})
	}
}

func TestListUsers(t *testing.T) {
	users := make([]db.User, 20)

	testCases := []testCase{
		{
			name:    "OK",
			payload: "0",
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().ListUsers(gomock.Any(), gomock.Eq(int32(0))).Times(1).Return(users, nil)
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
			name:    "OK",
			payload: "1",
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().ListUsers(gomock.Any(), gomock.Eq(int32(10))).Times(1).Return(users, nil)
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
			name:    "Invalid page",
			payload: "-1",
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().ListUsers(gomock.Any(), gomock.Eq(int32(0))).Times(0)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 400, rec.Code)
			},
		},
		{
			name:    "Missing query param defaults to page 0",
			payload: "",
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().ListUsers(gomock.Any(), gomock.Eq(int32(0))).Times(1).Return(users, nil)
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
			url := fmt.Sprintf("/api/v1/users?page=%s", tc.payload)

			req := httptest.NewRequest(http.MethodGet, url, nil)
			token, _, err := server.tokenMaker.CreateToken(uuid.New(), common.RandomString(12), server.tokenCfg.AccessTokenDuration())
			require.NoError(t, err)
			req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))

			server.router.ServeHTTP(rec, req)
			tc.checkResponse(rec)
		})
	}
}

func TestGetUserById(t *testing.T) {
	_, user := RandomUser(t)

	testCases := []testCase{
		{
			name:    "OK",
			payload: user.ID.String(),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetUserById(gomock.Any(), gomock.Eq(user.ID)).Times(1).Return(user, nil)
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
			name:    "Invalid ID",
			payload: common.RandomString(12),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetUserById(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 400, rec.Code)
			},
		},
		{
			name:    "NOT FOUND",
			payload: user.ID.String(),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetUserById(gomock.Any(), gomock.Any()).Times(1).Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 400, rec.Code)

				resp := new(response)

				body, err := io.ReadAll(rec.Body)
				require.NoError(t, err)

				require.NoError(t, json.Unmarshal(body, &resp))

				require.Equal(t, NOT_FOUND.Err, resp.Err)
			},
		},
		{
			name:    "INTERNAL ERROR",
			payload: user.ID.String(),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetUserById(gomock.Any(), gomock.Any()).Times(1).Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 500, rec.Code)

				resp := new(response)

				body, err := io.ReadAll(rec.Body)
				require.NoError(t, err)

				require.NoError(t, json.Unmarshal(body, &resp))
				require.NotNil(t, resp.Err)
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
			url := fmt.Sprintf("/api/v1/users/%s", tc.payload)

			req := httptest.NewRequest(http.MethodGet, url, nil)
			token, _, err := server.tokenMaker.CreateToken(user.ID, user.Username, server.tokenCfg.AccessTokenDuration())
			require.NoError(t, err)
			req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))

			server.router.ServeHTTP(rec, req)
			tc.checkResponse(rec)
		})
	}
}

func TestGetUserByUsername(t *testing.T) {
	_, user := RandomUser(t)

	testCases := []testCase{
		{
			name:    "OK",
			payload: user.Username,
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).Times(1).Return(user, nil)
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
			name:    "NOT FOUND",
			payload: user.Username,
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).Times(1).Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 400, rec.Code)

				resp := new(response)

				body, err := io.ReadAll(rec.Body)
				require.NoError(t, err)

				require.NoError(t, json.Unmarshal(body, &resp))

				require.Equal(t, NOT_FOUND.Err, resp.Err)
			},
		},
		{
			name:    "INTERNAL ERROR",
			payload: user.Username,
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).Times(1).Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 500, rec.Code)

				resp := new(response)

				body, err := io.ReadAll(rec.Body)
				require.NoError(t, err)

				require.NoError(t, json.Unmarshal(body, &resp))
				require.NotNil(t, resp.Err)
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
			url := fmt.Sprintf("/api/v1/users/one/%s", tc.payload)

			req := httptest.NewRequest(http.MethodGet, url, nil)

			server.router.ServeHTTP(rec, req)
			tc.checkResponse(rec)
		})
	}
}

type eqUpdateUserParamsMatcher struct {
	arg      db.UpdateUserPasswordParams
	password string
}

func (e *eqUpdateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.UpdateUserPasswordParams)
	if !ok {
		return false
	}

	err := common.CheckPassword([]byte(arg.Password), []byte(e.password))
	if err != nil {
		return false
	}

	e.arg.ID = arg.ID
	e.arg.Password = arg.Password
	e.arg.UpdatedAt = arg.UpdatedAt

	return reflect.DeepEqual(e.arg, arg)
}

func (e *eqUpdateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and id %v", e.arg, e.password)
}

func EqUpdateUserParams(arg *db.UpdateUserPasswordParams, password string) gomock.Matcher {
	return &eqUpdateUserParamsMatcher{*arg, password}
}

type eqUpdateUsernameParamsMatcher struct {
	arg      db.UpdateUsernameParams
	password string
}

func (e *eqUpdateUsernameParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.UpdateUsernameParams)
	if !ok {
		return false
	}

	e.arg.ID = arg.ID
	e.arg.Username = arg.Username
	e.arg.UpdatedAt = arg.UpdatedAt

	return reflect.DeepEqual(e.arg, arg)
}

func (e *eqUpdateUsernameParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and id %v", e.arg, e.password)
}

func EqUpdateUsernameParams(arg *db.UpdateUsernameParams, password string) gomock.Matcher {
	return &eqUpdateUsernameParamsMatcher{*arg, password}
}

func TestUpdateUser(t *testing.T) {
	_, user := RandomUser(t)
	session_id := uuid.New()
	newUsername := common.RandomString(12)
	newPassword := common.RandomString(12)

	testCases := []testCase{
		{
			name:    "OK-Username",
			payload: fmt.Sprintf(`{"field": "username", "payload": %q, "session_id": %q}`, newUsername, session_id),
			buildStubs: func(store *mock.MockStore) {
				arg := db.UpdateUsernameParams{
					ID:        user.ID,
					Username:  newUsername,
					UpdatedAt: user.UpdatedAt,
				}

				store.EXPECT().GetSessionById(gomock.Any(), gomock.Any()).Times(1)
				store.EXPECT().UpdateUsername(gomock.Any(), EqUpdateUsernameParams(&arg, newPassword)).Times(1).Return(user.ID, nil)
				store.EXPECT().DeleteSession(gomock.Any(), gomock.Any()).Times(1)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 200, rec.Code)

				resp := new(response)

				body, err := io.ReadAll(rec.Body)
				require.NoError(t, err)

				t.Log(resp.Err)
				require.NoError(t, json.Unmarshal(body, &resp))
				require.NotNil(t, resp.Data)
			},
		},
		{
			name:    "OK-Password",
			payload: fmt.Sprintf(`{"field": "password", "payload": %q, "session_id": %q}`, newPassword, session_id),
			buildStubs: func(store *mock.MockStore) {
				hashedPass, err := common.HashPassword(newPassword)
				require.NoError(t, err)

				arg := db.UpdateUserPasswordParams{
					Password:  hashedPass,
					ID:        user.ID,
					UpdatedAt: user.UpdatedAt,
				}
				store.EXPECT().GetSessionById(gomock.Any(), gomock.Any()).Times(1)
				store.EXPECT().UpdateUserPassword(gomock.Any(), EqUpdateUserParams(&arg, newPassword)).Times(1).Return(user.ID, nil)
				store.EXPECT().DeleteSession(gomock.Any(), gomock.Any()).Times(1)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 200, rec.Code)

				resp := new(response)

				body, err := io.ReadAll(rec.Body)
				require.NoError(t, err)

				t.Log(resp.Err)
				require.NoError(t, json.Unmarshal(body, &resp))
				require.NotNil(t, resp.Data)
			},
		},
		{
			name:    "Missing session ID",
			payload: fmt.Sprintf(`{"field": "username", "payload": %q}`, newUsername),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetSessionById(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 400, rec.Code)
			},
		},
		{
			name:    "Expired session",
			payload: fmt.Sprintf(`{"field": "username", "payload": %q, "session_id": %q}`, newUsername, session_id),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetSessionById(gomock.Any(), gomock.Eq(session_id)).Times(1).Return(db.Session{}, sql.ErrNoRows)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 401, rec.Code)
			},
		},
		{
			name:    "Unknown field",
			payload: fmt.Sprintf(`{"field": "unknown", "payload": %q, "session_id": %q}`, newUsername, session_id),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetSessionById(gomock.Any(), gomock.Eq(session_id)).Times(0)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 400, rec.Code)
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

			url := fmt.Sprintf("/api/v1/users/%s", user.ID)

			req := httptest.NewRequest(http.MethodPatch, url, strings.NewReader(tc.payload))

			token, _, err := server.tokenMaker.CreateToken(user.ID, user.Username, server.tokenCfg.AccessTokenDuration())
			require.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))

			rec := httptest.NewRecorder()

			server.router.ServeHTTP(rec, req)
			tc.checkResponse(rec)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	_, user := RandomUser(t)
	sessionId := uuid.New()

	testCases := []testCase{
		{
			name:    "OK",
			payload: fmt.Sprintf(`{"session_id": %q}`, sessionId),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetSessionById(gomock.Any(), gomock.Eq(sessionId)).Times(1)
				store.EXPECT().DeleteOneUser(gomock.Any(), gomock.Eq(user.ID)).Times(1).Return(user.ID, nil)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 200, rec.Code)
				resp := new(response)

				body, err := io.ReadAll(rec.Body)
				require.NoError(t, err)
				require.NoError(t, json.Unmarshal(body, &resp))
				require.NotNil(t, resp.Data)
				require.Empty(t, resp.Err)
			},
		},
		{
			name:    "Invalid session",
			payload: fmt.Sprintf(`{"session_id": %q}`, sessionId),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetSessionById(gomock.Any(), gomock.Eq(sessionId)).Times(1).Return(db.Session{}, sql.ErrNoRows)

			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 401, rec.Code)
				resp := new(response)

				body, err := io.ReadAll(rec.Body)
				require.NoError(t, err)
				require.NoError(t, json.Unmarshal(body, &resp))
				require.NotNil(t, resp.Err)
				require.Empty(t, resp.Data)
			},
		},
		{
			name:    "Internal error",
			payload: fmt.Sprintf(`{"session_id": %q}`, sessionId),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetSessionById(gomock.Any(), gomock.Eq(sessionId)).Times(1).Return(db.Session{}, sql.ErrConnDone)

			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 500, rec.Code)
				resp := new(response)

				body, err := io.ReadAll(rec.Body)
				require.NoError(t, err)
				require.NoError(t, json.Unmarshal(body, &resp))
				require.NotNil(t, resp.Err)
				require.Empty(t, resp.Data)
			},
		},
		{
			name:    "Missing body payload",
			payload: "",
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetSessionById(gomock.Any(), gomock.Eq(sessionId)).Times(0)

			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 400, rec.Code)

				resp := new(response)

				body, err := io.ReadAll(rec.Body)
				require.NoError(t, err)

				require.NoError(t, json.Unmarshal(body, &resp))
				require.Contains(t, resp.Err, "SessionId", "required")
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

			url := fmt.Sprintf("/api/v1/users/%s", user.ID)

			req := httptest.NewRequest(http.MethodDelete, url, strings.NewReader(tc.payload))

			token, _, err := server.tokenMaker.CreateToken(user.ID, user.Username, server.tokenCfg.AccessTokenDuration())
			require.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))

			rec := httptest.NewRecorder()

			server.router.ServeHTTP(rec, req)
			tc.checkResponse(rec)
		})
	}
}
