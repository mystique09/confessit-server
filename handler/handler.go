// Package specification CNFS API.
//
// # Documentation for the CNFS API.
//
// Schemes: http
// BasePath: /api/v1
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
//	securityDefinitions:
//		key:
//		 type: apiKey
//		 in: header
//		 name: authorization
//
// swagger:meta
package handler

import (
	"cnfs/common"
	db "cnfs/db/sqlc"
	"cnfs/token"
	"cnfs/web"
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
		cfg        *config.Config
		router     *echo.Echo
		store      db.Store
		tokenMaker token.Maker
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
	conn := common.SetupDb(cfg.DatabaseUrl)
	store := db.NewStore(conn, cfg)
	server, err := NewServer(store, cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Fatal(server.router.Start(cfg.Host + ":" + cfg.Port))
}

func NewServer(store db.Store, cfg *config.Config) (*Server, error) {

	tokenMaker, err := token.NewPasetoMaker(cfg.PasetoSymmetricKey)
	if err != nil {
		log.Fatalf("cannot make token maker: %s", err.Error())
	}

	server := &Server{
		cfg:        cfg,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()
	return server, nil
}

func (s *Server) setupRouter() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	logger := zerolog.New(os.Stdout)

	e.Use(s.loggerMiddleware(&logger))
	e.Use(s.corsMiddleware())
	e.Use(s.rateLimitMiddleware(20))

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "CNFS server, running on Echo v4.")
	})
	e.GET("/health", func(c echo.Context) error {
		return c.String(200, "Health: 100/100 Armor: 100/100")
	})
	e.StaticFS("/docs", web.BuildHttpFS())

	auth := e.Group("/api/v1/auth")
	auth.POST("", s.loginUser)
	auth.POST("/refresh", s.refreshAccessToken)
	auth.POST("/validate", s.validateAccessToken)
	auth.DELETE("/clear", s.logoutUser)

	users := e.Group("/api/v1/users")
	users.GET("", s.listUsers, s.authMiddleware)
	users.POST("", s.createUser)
	users.GET("/:id", s.getUserById, s.authMiddleware)
	users.GET("/:id/messages", s.listMessages, s.authMiddleware)
	users.PATCH("/:id", s.updateUser, s.authMiddleware)
	users.DELETE("/:id", s.deleteUser, s.authMiddleware)
	users.GET("/one/:username", s.getUserByUsername)

	messages := e.Group("/api/v1/messages")
	messages.GET("/:id", s.getMessageById, s.authMiddleware)
	messages.POST("", s.createMessage)
	messages.PUT("/:id", s.updateMessage, s.authMiddleware)
	messages.DELETE("/:id", s.deleteMessage, s.authMiddleware)

	posts := e.Group("/api/v1/posts")
	posts.GET("", s.listAllPosts)
	posts.GET("/:id", s.getPostById)
	posts.POST("", s.createNewPost, s.authMiddleware)
	posts.PATCH("/:id", s.updatePost, s.authMiddleware)
	posts.DELETE("/:id", s.deletePost, s.authMiddleware)
	posts.GET("/:id/comments", s.listAllComments)

	comments := e.Group("/api/v1/comments")
	comments.GET("/:id", s.getCommentById)
	comments.POST("", s.createComment, s.authMiddleware)
	comments.PUT("/:id", s.updateComment, s.authMiddleware)
	comments.DELETE("/:id", s.deleteComment, s.authMiddleware)

	s.router = e
}
