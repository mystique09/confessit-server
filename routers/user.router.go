package routers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (r *Route) GetUsers(c echo.Context) error {
	return c.String(http.StatusOK, "Get users")
}

func (r *Route) GetUser(c echo.Context) error {
	return c.String(http.StatusOK, "Get user")
}

func (r *Route) UpdateUser(c echo.Context) error {
	return c.String(http.StatusOK, "Update user")
}

func (r *Route) DeleteUser(c echo.Context) error {
	return c.String(http.StatusOK, "Delete user")
}
