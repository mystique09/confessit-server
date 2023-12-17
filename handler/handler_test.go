package handler

import (
	"cnfs/common"
	"cnfs/config"
	"cnfs/db/mock"
	db "cnfs/db/sqlc"
	"cnfs/domain"
	"log"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	name          string
	payload       string
	buildStubs    func(store *mock.MockStore)
	checkResponse func(rec *httptest.ResponseRecorder)
}

var cfg *domain.IConfig

func TestMain(m *testing.M) {
	configLoader := config.NewConfigLoader("..", "app")
	serverConfig, err := config.NewServerConfig(configLoader)
	if err != nil {
		log.Fatal(err)
	}
	tokenConfig, err := config.NewTokenConfig(configLoader)
	if err != nil {
		log.Fatal(err)

	}
	c := config.NewConfig(serverConfig, tokenConfig)

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
		ID:        uuid.New(),
		Username:  common.RandomString(12),
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return password, user
}

func RandomMessage(t *testing.T, userId uuid.UUID) db.Message {
	return db.Message{
		ID:         uuid.New(),
		ReceiverID: userId,
		Content:    common.RandomString(48),
		Seen:       false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func RandomPost(t *testing.T, userIdentityId uuid.UUID) db.Post {
	return db.Post{
		ID:             uuid.New(),
		UserIdentityID: userIdentityId,
		Content:        common.RandomString(48),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

func RandomComment(t *testing.T, postId, parentId uuid.UUID) db.Comment {
	return db.Comment{
		ID:             uuid.New(),
		PostID:         postId,
		Content:        common.RandomString(48),
		UserIdentityID: uuid.New(),
		ParentID:       parentId,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}
