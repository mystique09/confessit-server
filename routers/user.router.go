package routers

import (
	"confessit/handlers"
	"confessit/models"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (r *Route) GetUsers(c echo.Context) error {
	var users []models.UserResponse = handlers.GetUsers(r.Conn)

	return c.JSON(http.StatusOK, users)
}

func (r *Route) GetUser(c echo.Context) error {
	var unameParam string = c.Param("name")

	if unameParam == "" {
		return c.String(http.StatusBadRequest, "missing name parameter field")
	}

	var user models.User = handlers.GetUser(r.Conn, unameParam)

	if user.Username == "" || user.ID == uuid.Nil {
		return c.String(http.StatusBadRequest, "user does not exist.")
	}

	return c.JSON(http.StatusOK, user.ToResponse())
}

func (r *Route) UpdateUser(c echo.Context) error {
	uidParam, uidError := uuid.Parse(c.Param("name"))
	if uidError != nil {
		return c.JSON(http.StatusBadRequest, uidError.Error())
	}

	var payload models.UserUpdatePayload

	if err := (&echo.DefaultBinder{}).BindBody(c, &payload); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := validate.Struct(payload); err != nil {
		return c.JSON(http.StatusOK, err.Error())
	}

	hasUser := handlers.GetUserById(r.Conn, uidParam)

	if hasUser.Username == "" || hasUser.ID == uuid.Nil {
		return c.JSON(http.StatusBadRequest, "user does not exist.")
	}

	checkUser := handlers.GetUser(r.Conn, payload.Username)

	if checkUser.Username != "" || checkUser.ID != uuid.Nil {
		return c.JSON(http.StatusBadRequest, "username already taken.")
	}

	if err := handlers.UpdateUser(r.Conn, uidParam, payload); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, "user updated.")
}

func (r *Route) DeleteUser(c echo.Context) error {
	uidParam, uidError := uuid.Parse(c.Param("name"))

	if uidError != nil {
		return c.String(http.StatusBadRequest, uidError.Error())
	}

	hasUser := handlers.GetUserById(r.Conn, uidParam)

	if hasUser.Username == "" || hasUser.ID == uuid.Nil {
		return c.String(http.StatusBadRequest, "user does not exist.")
	}

	if err := handlers.DeleteUser(r.Conn, uidParam); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, "user deleted")
}
