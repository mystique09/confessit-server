package routers

import (
	"confessit/handlers"
	"confessit/models"
	"confessit/utils"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func (r *Route) GetMessages(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	username := utils.GetPayloadUsername(user)

	messages := handlers.GetMessages(r.Conn, username)
	return c.JSON(http.StatusOK, NewResponse("All messages", messages))
}

func (r *Route) GetMessage(c echo.Context) error {
	return c.String(http.StatusOK, "Get message")
}

func (r *Route) CreateMessage(c echo.Context) error {
	var payload models.MessageCreatePayload

	if err := (&echo.DefaultBinder{}).BindBody(c, &payload); err != nil {
		return c.JSON(http.StatusBadRequest, NewError(err.Error()))
	}

	if err := handlers.CreateMessage(r.Conn, payload); err != nil {
		return c.JSON(http.StatusBadRequest, NewError(err.Error()))
	}

	return c.JSON(http.StatusOK, NewResponse("Create message", payload))
}

func (r *Route) DeleteMessage(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	username := utils.GetPayloadUsername(user)

	var payload models.MessageDeletePayload

	if err := (&echo.DefaultBinder{}).BindBody(c, &payload); err != nil {
		return c.JSON(http.StatusBadRequest, NewError(err.Error()))
	}

	if err := handlers.DeleteMessage(r.Conn, payload.ID, username); err != nil {
		return c.JSON(http.StatusBadRequest, NewError(err.Error()))
	}

	return c.JSON(http.StatusOK, NewResponse("1 message deleted", payload))
}
