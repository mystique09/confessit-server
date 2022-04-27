package routers

import (
	"confessit/handlers"
	"confessit/models"
	"confessit/utils"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var NewError = utils.NewError
var NewResponse = utils.NewResponse

func (r *Route) GetUsers(c echo.Context) error {
	var users []models.UserResponse = handlers.GetUsers(r.Conn)

	return c.JSON(http.StatusOK, NewResponse("All users", users))
}

func (r *Route) GetUser(c echo.Context) error {
	var unameParam string = c.Param("name")

	if unameParam == "" {
		return c.JSON(http.StatusBadRequest, NewError("missing name parameter field"))
	}

	var user models.User = handlers.GetUser(r.Conn, unameParam)

	if user.Username == "" || user.ID == uuid.Nil {
		return c.JSON(http.StatusBadRequest, NewError("user does not exist."))
	}

	return c.JSON(http.StatusOK, NewResponse("One user", user.ToResponse()))
}

func (r *Route) GetUserById(c echo.Context) error {
	user_token := c.Get("user").(*jwt.Token)
	username := utils.GetPayloadUsername(user_token)

	user := handlers.GetUser(r.Conn, username)

	return c.JSON(http.StatusOK, NewResponse("One user", user.ToResponse()))
}

func (r *Route) UpdateUser(c echo.Context) error {
	uidParam, uidError := uuid.Parse(c.Param("name"))
	if uidError != nil {
		return c.JSON(http.StatusBadRequest, NewError(uidError.Error()))
	}

	var payload models.UserUpdatePayload

	if err := (&echo.DefaultBinder{}).BindBody(c, &payload); err != nil {
		return c.JSON(http.StatusBadRequest, NewError(err.Error()))
	}

	if err := validate.Struct(payload); err != nil {
		return c.JSON(http.StatusOK, err.Error())
	}

	hasUser := handlers.GetUserById(r.Conn, uidParam)

	if hasUser.Username == "" || hasUser.ID == uuid.Nil {
		return c.JSON(http.StatusBadRequest, NewError("user does not exist."))
	}

	checkUser := handlers.GetUser(r.Conn, payload.Username)

	if checkUser.Username != "" || checkUser.ID != uuid.Nil {
		return c.JSON(http.StatusBadRequest, NewError("username already taken."))
	}

	if err := handlers.UpdateUser(r.Conn, uidParam, payload); err != nil {
		return c.JSON(http.StatusBadRequest, NewError(err.Error()))
	}

	return c.JSON(http.StatusOK, NewResponse("user updated.", uidParam))
}

func (r *Route) DeleteUser(c echo.Context) error {
	uidParam, uidError := uuid.Parse(c.Param("name"))

	if uidError != nil {
		return c.JSON(http.StatusBadRequest, NewError(uidError.Error()))
	}

	hasUser := handlers.GetUserById(r.Conn, uidParam)

	if hasUser.Username == "" || hasUser.ID == uuid.Nil {
		return c.JSON(http.StatusBadRequest, NewError("user does not exist."))
	}

	if err := handlers.DeleteUser(r.Conn, uidParam); err != nil {
		return c.JSON(http.StatusBadRequest, NewError(err.Error()))
	}

	return c.JSON(http.StatusOK, NewResponse("user deleted", uidParam))
}
