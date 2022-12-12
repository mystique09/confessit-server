package handler

import (
	"cnfs/db/mock"
	db "cnfs/db/sqlc"
	"cnfs/utils"
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

	err := utils.CheckPassword([]byte(arg.Password), []byte(e.password))
	if err != nil {
		return false
	}

	e.arg.ID = arg.ID
	e.arg.Password = arg.Password

	return reflect.DeepEqual(e.arg, arg)
}

func (e *eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and id %v", e.arg, e.password)
}

func EqCreateUserParams(arg *db.CreateUserParams, password string) gomock.Matcher {
	return &eqCreateUserParamsMatcher{*arg, password}
}

func TestCreateUser(t *testing.T) {
	password, user := randomUser(t)

	testCases := []testCase{
		{
			name:    "OK",
			payload: fmt.Sprintf(`{"username": %q, "password": %q}`, user.Username, password),
			buildStubs: func(store *mock.MockStore) {
				arg := db.CreateUserParams{
					ID:       user.ID,
					Username: user.Username,
					Password: user.Password,
				}

				store.
					EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(&arg, password)).
					Times(1).
					Return(user.ID, nil)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 200, rec.Code)
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
				arg := db.CreateUserParams{
					ID:       user.ID,
					Username: user.Username,
					Password: user.Password,
				}

				store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(&arg, password)).Times(1).Return(user.ID, errors.New("unique violation"))
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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock.NewMockStore(ctrl)
			tc.buildStubs(store)

			server, err := NewServer(store, cfg)
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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock.NewMockStore(ctrl)
			tc.buildStubs(store)

			server, err := NewServer(store, cfg)
			require.NoError(t, err)

			rec := httptest.NewRecorder()
			url := fmt.Sprintf("/api/v1/users?page=%s", tc.payload)

			req := httptest.NewRequest(http.MethodGet, url, nil)

			server.router.ServeHTTP(rec, req)
			tc.checkResponse(rec)
		})
	}
}

func TestGetUserById(t *testing.T) {
	_, user := randomUser(t)

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
			payload: utils.RandomString(12),
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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock.NewMockStore(ctrl)
			tc.buildStubs(store)

			server, err := NewServer(store, cfg)
			require.NoError(t, err)

			rec := httptest.NewRecorder()
			url := fmt.Sprintf("/api/v1/users/%s", tc.payload)

			req := httptest.NewRequest(http.MethodGet, url, nil)

			server.router.ServeHTTP(rec, req)
			tc.checkResponse(rec)
		})
	}
}
