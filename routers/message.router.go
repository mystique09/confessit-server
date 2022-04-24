package routers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (r *Route) GetMessages(c echo.Context) error {
	return c.String(http.StatusOK, "Get messages")
}

func (r *Route) GetMessage(c echo.Context) error {
	return c.String(http.StatusOK, "Get message")
}

func (r *Route) CreateMessage(c echo.Context) error {
	return c.String(http.StatusOK, "Create message")
}

func (r *Route) DeleteMessage(c echo.Context) error {
	return c.String(http.StatusOK, "Delete messages")
}
