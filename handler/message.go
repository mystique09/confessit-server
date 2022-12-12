package handler

import "github.com/labstack/echo/v4"

func (s *Server) createMessage(c echo.Context) error {
	return c.JSON(200, "create Message")
}

func (s *Server) listMessages(c echo.Context) error {
	return c.JSON(200, "all Messages")
}

func (s *Server) getMessageById(c echo.Context) error {
	return c.JSON(200, "Message by id")
}

func (s *Server) deleteMessage(c echo.Context) error {
	return c.JSON(200, "delete Message")
}
