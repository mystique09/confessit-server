package handler

import (
	"cnfs/db/mock"
	db "cnfs/db/sqlc"
	"cnfs/utils"
	"database/sql"
	"encoding/json"
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

type eqCreateMessageParams struct {
	arg db.CreateMessageParams
	id  uuid.UUID
}

func (e *eqCreateMessageParams) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateMessageParams)
	if !ok {
		return false
	}

	e.arg.ID = arg.ID

	return reflect.DeepEqual(e.arg, arg)
}

func (e *eqCreateMessageParams) String() string {
	return fmt.Sprintf("matches arg %v and id %v", e.arg, e.id)
}

func EqCreateMessageParams(arg *db.CreateMessageParams, id uuid.UUID) gomock.Matcher {
	return &eqCreateMessageParams{*arg, id}
}

func TestServer_createMessage(t *testing.T) {
	_, user := randomUser(t)
	msg := randomMessage(t, user.ID)

	testCases := []testCase{
		{
			name:    "OK",
			payload: fmt.Sprintf(`{"receiver_id": %q, "content": %q}`, user.ID, msg.Content),
			buildStubs: func(store *mock.MockStore) {
				arg := db.CreateMessageParams{
					ID:         msg.ID,
					ReceiverID: user.ID,
					Content:    msg.Content,
				}
				store.EXPECT().GetUserById(gomock.Any(), gomock.Eq(user.ID)).Times(1).Return(user, nil)
				store.EXPECT().CreateMessage(gomock.Any(), EqCreateMessageParams(&arg, arg.ID)).Times(1).Return(msg, nil)
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
			name:    "USER NOT FOUND",
			payload: fmt.Sprintf(`{"receiver_id": %q, "content": %q}`, user.ID, msg.Content),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetUserById(gomock.Any(), gomock.Eq(user.ID)).Times(1).Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 404, rec.Code)
				resp := new(response)

				body, err := io.ReadAll(rec.Body)
				require.NoError(t, err)
				require.NoError(t, json.Unmarshal(body, &resp))
				require.NotNil(t, resp.Err)
			},
		},
		{
			name:    "INTERNAL ERROR",
			payload: fmt.Sprintf(`{"receiver_id": %q, "content": %q}`, user.ID, msg.Content),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetUserById(gomock.Any(), gomock.Eq(user.ID)).Times(1).Return(db.User{}, sql.ErrConnDone)
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
		{
			name:    "INVALID PAYLOAD/Missing receiver ID",
			payload: fmt.Sprintf(`{"content": %q}`, msg.Content),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetUserById(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 400, rec.Code)
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

			req := httptest.NewRequest(http.MethodPost, "/api/v1/messages", strings.NewReader(tc.payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()

			server.router.ServeHTTP(rec, req)
			tc.checkResponse(rec)
		})
	}
}

func TestListMessages(t *testing.T) {
	_, user := randomUser(t)
	messages := make([]db.Message, 20)
	testCases := []testCase{
		{
			name:    "OK - w/ page 0",
			payload: fmt.Sprintf("/api/v1/users/%s/messages?page=%d", user.ID, 0),
			buildStubs: func(store *mock.MockStore) {
				arg := db.ListMessageParams{
					ReceiverID: user.ID,
					Offset:     0,
				}
				store.EXPECT().ListMessage(gomock.Any(), gomock.Eq(arg)).Times(1).Return(messages, nil)
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
			name:    "OK - w/ empty page query",
			payload: fmt.Sprintf("/api/v1/users/%s/messages", user.ID),
			buildStubs: func(store *mock.MockStore) {
				arg := db.ListMessageParams{
					ReceiverID: user.ID,
					Offset:     0,
				}
				store.EXPECT().ListMessage(gomock.Any(), gomock.Eq(arg)).Times(1).Return(messages, nil)
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
			name:    "OK - w/ page query > 0",
			payload: fmt.Sprintf("/api/v1/users/%s/messages?page=%d", user.ID, 1),
			buildStubs: func(store *mock.MockStore) {
				arg := db.ListMessageParams{
					ReceiverID: user.ID,
					Offset:     10,
				}
				store.EXPECT().ListMessage(gomock.Any(), gomock.Eq(arg)).Times(1).Return(messages, nil)
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
			name:    "User NOT FOUND",
			payload: fmt.Sprintf("/api/v1/users/%s/messages?page=%d", user.ID, 0),
			buildStubs: func(store *mock.MockStore) {
				arg := db.ListMessageParams{
					ReceiverID: user.ID,
					Offset:     0,
				}
				store.EXPECT().ListMessage(gomock.Any(), gomock.Eq(arg)).Times(1).Return([]db.Message{}, sql.ErrNoRows)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 404, rec.Code)
				resp := new(response)

				body, err := io.ReadAll(rec.Body)
				require.NoError(t, err)
				require.NoError(t, json.Unmarshal(body, &resp))
				require.NotNil(t, resp.Err)
				require.Empty(t, resp.Data)
			},
		},
		{
			name:    "NEGATIVE page query not allowed",
			payload: fmt.Sprintf("/api/v1/users/%s/messages?page=%d", user.ID, -1),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().ListMessage(gomock.Any(), gomock.Any()).Times(0)
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

			req := httptest.NewRequest(http.MethodGet, tc.payload, nil)
			token, _, err := server.tokenMaker.CreateToken(uuid.New(), utils.RandomString(12), cfg.AccessTokenDuration)
			require.NoError(t, err)
			req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))

			rec := httptest.NewRecorder()

			server.router.ServeHTTP(rec, req)
			tc.checkResponse(rec)
		})
	}
}

func TestGetMessageById(t *testing.T) {
	_, user := randomUser(t)
	msg := randomMessage(t, user.ID)

	testCases := []testCase{
		{
			name:    "OK",
			payload: fmt.Sprintf("/api/v1/messages/%s", msg.ID),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetMessageById(gomock.Any(), gomock.Eq(msg.ID)).Times(1).Return(msg, nil)
				store.EXPECT().UpdateMessageStatus(gomock.Any(), gomock.Any()).Times(1)
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
			name:    "NOT FOUND",
			payload: fmt.Sprintf("/api/v1/messages/%s", msg.ID),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetMessageById(gomock.Any(), gomock.Eq(msg.ID)).Times(1).Return(db.Message{}, sql.ErrNoRows)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 404, rec.Code)
				resp := new(response)

				body, err := io.ReadAll(rec.Body)
				require.NoError(t, err)
				require.NoError(t, json.Unmarshal(body, &resp))
				require.NotNil(t, resp.Err)
				require.Empty(t, resp.Data)
			},
		},
		{
			name:    "INTERNAL ERROR",
			payload: fmt.Sprintf("/api/v1/messages/%s", msg.ID),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetMessageById(gomock.Any(), gomock.Eq(msg.ID)).Times(1).Return(db.Message{}, sql.ErrConnDone)
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
			name:    "Invalid UUID",
			payload: fmt.Sprintf("/api/v1/messages/%s", utils.RandomString(12)),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetMessageById(gomock.Any(), gomock.Any()).Times(0)
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

			req := httptest.NewRequest(http.MethodGet, tc.payload, nil)
			token, _, err := server.tokenMaker.CreateToken(uuid.New(), utils.RandomString(12), cfg.AccessTokenDuration)
			require.NoError(t, err)
			req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))

			rec := httptest.NewRecorder()

			server.router.ServeHTTP(rec, req)
			tc.checkResponse(rec)
		})
	}
}

func TestDeleteMessage(t *testing.T) {
	_, user := randomUser(t)
	msg := randomMessage(t, user.ID)

	testCases := []testCase{
		{
			name:    "OK",
			payload: fmt.Sprintf("/api/v1/messages/%s", msg.ID),
			buildStubs: func(store *mock.MockStore) {
				arg := db.DeleteOneMessageParams{
					ReceiverID: msg.ReceiverID,
					ID:         msg.ID,
				}
				store.EXPECT().GetMessageById(gomock.Any(), gomock.Eq(msg.ID)).Times(1).Return(msg, nil)
				store.EXPECT().DeleteOneMessage(gomock.Any(), gomock.Eq(arg)).Times(1).Return(msg.ID, nil)
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
			name:    "NOT FOUND",
			payload: fmt.Sprintf("/api/v1/messages/%s", msg.ID),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetMessageById(gomock.Any(), gomock.Eq(msg.ID)).Times(1).Return(db.Message{}, sql.ErrNoRows)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, 404, rec.Code)
				resp := new(response)

				body, err := io.ReadAll(rec.Body)
				require.NoError(t, err)
				require.NoError(t, json.Unmarshal(body, &resp))
				require.NotNil(t, resp.Err)
				require.Empty(t, resp.Data)
			},
		},
		{
			name:    "INTERNAL ERROR",
			payload: fmt.Sprintf("/api/v1/messages/%s", msg.ID),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetMessageById(gomock.Any(), gomock.Eq(msg.ID)).Times(1).Return(db.Message{}, sql.ErrConnDone)
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
			name:    "UNAUTHORIZED ACCESS",
			payload: fmt.Sprintf("/api/v1/messages/%s", msg.ID),
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().GetMessageById(gomock.Any(), gomock.Eq(msg.ID)).Times(1).Return(randomMessage(t, uuid.New()), nil)
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

			req := httptest.NewRequest(http.MethodDelete, tc.payload, nil)
			token, _, err := server.tokenMaker.CreateToken(user.ID, user.Username, cfg.AccessTokenDuration)
			require.NoError(t, err)
			req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))

			rec := httptest.NewRecorder()

			server.router.ServeHTTP(rec, req)
			tc.checkResponse(rec)
		})
	}
}
