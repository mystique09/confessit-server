package handler

import (
	db "cnfs/db/sqlc"
	"cnfs/utils"
	"log"
	"net/http"
	"os"

	"cnfs/config"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type (
	Server struct {
		cfg    *config.Config
		router *echo.Echo
		store  db.Store
	}

	CustomValidator struct {
		validator *validator.Validate
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func Launch(cfg *config.Config) {
	conn := utils.SetupDb(cfg.DatabaseUrl)
	store := db.NewStore(conn, cfg)
	server, err := NewServer(&store, cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Fatal(server.router.Start(cfg.Host))
}

func NewServer(store *db.Store, cfg *config.Config) (*Server, error) {
	server := &Server{
		cfg:   cfg,
		store: *store,
	}

	server.setupRouter()

	return server, nil
}

func (s *Server) setupRouter() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	logger := zerolog.New(os.Stdout)

	e.Use(loggerMiddleware(&logger))

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "CNFS server, running on Echo v4.")
	})

	e.GET("/health", func(c echo.Context) error {
		return c.String(200, "Health: 100/100 Armor: 100/100")
	})

	auth := e.Group("/api/v1/auth")
	auth.POST("", s.loginUser)
	auth.POST("/refresh", s.refreshAccessToken)
	auth.POST("/validate", s.validateAccessToken)

	users := e.Group("/api/v1/users")
	users.GET("", s.listUsers)
	users.GET("/:id", s.getUserById)
	users.POST("", s.createUser)
	users.PUT("/:id", s.updateUser)
	users.DELETE("/:id", s.deleteUser)

	messages := e.Group("/api/v1/messages")
	messages.GET("", s.listMessages)
	messages.GET("/:id", s.getMessageById)
	messages.POST("", s.createMessage)
	messages.DELETE("/:id", s.deleteMessage)

	s.router = e
}
