package handler

import (
	"cnfs/config"
	"cnfs/db/mock"
	db "cnfs/db/sqlc"
	"cnfs/utils"
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

func randomUser(t *testing.T) (string, db.User) {
	password := utils.RandomString(14)
	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)

	user := db.User{
		ID:       uuid.New(),
		Username: utils.RandomString(12),
		Password: hashedPassword,
	}

	return password, user
}

func randomMessage(t *testing.T, userId uuid.UUID) db.Message {
	return db.Message{
		ID:         uuid.New(),
		ReceiverID: userId,
		Content:    utils.RandomString(48),
		Seen:       false,
	}
}
