package routers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Route struct {
	Conn *gorm.DB
}

func ConfessIt(c echo.Context) error {
	return c.String(http.StatusOK, "Confess anonymously!")
}
