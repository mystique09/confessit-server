package routers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func ConfessIt(c echo.Context) error {
	return c.String(http.StatusOK, "Confess anonymously!")
}
