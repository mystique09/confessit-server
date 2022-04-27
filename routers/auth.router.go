package routers

import (
	"confessit/handlers"
	"confessit/models"
	"confessit/utils"
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
		return c.JSON(http.StatusBadRequest, err.Error())
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
	var payload models.UserLoginPayload

	if err := (&echo.DefaultBinder{}).BindBody(c, &payload); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := validate.Struct(payload); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	hasUser := handlers.GetUser(r.Conn, payload.Username)

	if hasUser.Username == "" {
		return c.JSON(http.StatusUnauthorized, "user does not exist.")
	}

	if err := hasUser.ValidatePassword(payload.Password); err != nil {
		return c.JSON(http.StatusUnauthorized, "password mismatch.")
	}

	token, err := utils.CreateJwt(models.JwtUserPayload{
		Id:       hasUser.ID,
		Username: hasUser.Username,
	})

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, token)
}
