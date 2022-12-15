package web

import (
	"embed"
	"io/fs"

	"github.com/labstack/echo/v4"
)

//go:embed all:swagger
var SwaggerFS embed.FS

func BuildHttpFS() fs.FS {
	s := echo.MustSubFS(SwaggerFS, "swagger")
	return s
}
