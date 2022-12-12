package handler

import "github.com/labstack/echo/v4"

func (s *Server) loginUser(c echo.Context) error {
	return c.JSON(200, "logged in")
}

func (s *Server) refreshAccessToken(c echo.Context) error {
	return c.JSON(200, "new token is assigned")
}

func (s *Server) validateAccessToken(c echo.Context) error {
	return c.JSON(200, "current token is valid")
}
