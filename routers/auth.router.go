package routers

import (
	"confessit/handlers"
	"confessit/models"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var validate *validator.Validate = validator.New()

func (r *Route) Signup(c echo.Context) error {
	var payload models.UserCreatePayload

	if err := (&echo.DefaultBinder{}).BindBody(c, &payload); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := validate.Struct(payload); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	hasUser := handlers.GetUser(r.Conn, payload.Username)

	if hasUser.Username != "" {
		return c.JSON(http.StatusBadRequest, "user already exist.")

	}

	if err := handlers.CreateUser(r.Conn, payload); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, "New user created.")
}

func (r *Route) Login(c echo.Context) error {
	return c.String(http.StatusOK, "Login")
}
