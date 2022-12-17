package handler

import (
	"cnfs/common"
	"cnfs/config"
	"cnfs/db/mock"
	db "cnfs/db/sqlc"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	name          string
	payload       string
	buildStubs    func(store *mock.MockStore)
	checkResponse func(rec *httptest.ResponseRecorder)
}

var cfg *config.Config

func TestMain(m *testing.M) {
	c, err := config.LoadConfig("..", "app")
	if err != nil {
		log.Fatal(err.Error())
	}

	cfg = &c

	os.Exit(m.Run())
}

func RandomUser(t *testing.T) (string, db.User) {
	password := common.RandomString(14)
	hashedPassword, err := common.HashPassword(password)
	require.NoError(t, err)

	user := db.User{
		ID:       uuid.New(),
		Username: common.RandomString(12),
		Password: hashedPassword,
	}

	return password, user
}

func RandomMessage(t *testing.T, userId uuid.UUID) db.Message {
	return db.Message{
		ID:         uuid.New(),
		ReceiverID: userId,
		Content:    common.RandomString(48),
		Seen:       false,
	}
}

func RandomPost(t *testing.T, id uuid.UUID) db.Post {
	return db.Post{
		ID:             uuid.New(),
		UserIdentityID: id,
		Content:        common.RandomString(48),
	}
}
