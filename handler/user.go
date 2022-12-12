package handler

import "github.com/labstack/echo/v4"

func (s *Server) createUser(c echo.Context) error {
	return c.JSON(200, "create user")
}

func (s *Server) listUsers(c echo.Context) error {
	return c.JSON(200, "all Users")
}

func (s *Server) getUserById(c echo.Context) error {
	return c.JSON(200, "User by id")
}

func (s *Server) updateUser(c echo.Context) error {
	return c.JSON(200, "update User")
}

func (s *Server) deleteUser(c echo.Context) error {
	return c.JSON(200, "delete User")
}
