package routers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (r *Route) Signup(c echo.Context) error {
	return c.String(http.StatusOK, "Signup")
}

func (r *Route) Login(c echo.Context) error {
	return c.String(http.StatusOK, "Login")
}
