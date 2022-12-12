package handler

import (
	db "cnfs/db/sqlc"
	"cnfs/token"
	"cnfs/utils"
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type (
	createUserRequest struct {
		Username string `json:"username" validate:"gte=8,required"`
		Password string `json:"password" validate:"gte=8,required"`
	}

	updateUserRequest struct {
		Field   string `json:"field" validate:"gte=8,required"`
		Payload string `json:"payload" validate:"gte=8,required"`
	}
)

func (s *Server) createUser(c echo.Context) error {
	var data createUserRequest

	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

	if err := c.Validate(&data); err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

	hashedPassword, err := utils.HashPassword(data.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

	createUserParam := db.CreateUserParams{
		ID:       uuid.New(),
		Username: data.Username,
		Password: hashedPassword,
	}

	user, err := s.store.CreateUser(c.Request().Context(), createUserParam)
	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			return c.JSON(http.StatusBadRequest, newError("user already exist"))
		}
		return c.JSON(http.StatusInternalServerError, newError(err.Error()))
	}

	return c.JSON(http.StatusOK, newResponse(user))
}

func (s *Server) listUsers(c echo.Context) error {
	pageParam := c.QueryParam("page")
	if pageParam == "" {
		pageParam = "0"
	}

	page, err := strconv.ParseUint(pageParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

	users, err := s.store.ListUsers(c.Request().Context(), int32(page)*10)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

	return c.JSON(200, newResponse(users))
}

func (s *Server) getUserById(c echo.Context) error {
	id := c.Param("id")

	userId, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

	user, err := s.store.GetUserById(c.Request().Context(), userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusBadRequest, NOT_FOUND)
		}
		return c.JSON(http.StatusInternalServerError, newError(err.Error()))
	}

	return c.JSON(200, newResponse(user))
}

func (s *Server) updateUser(c echo.Context) error {
	id := c.Param("id")
	userId, err := uuid.Parse(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

	data := new(updateUserRequest)

	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

	if err := c.Validate(&data); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	tokenPayload, ok := c.Get("user").(*token.Payload)
	if !ok {
		return c.JSON(http.StatusUnauthorized, INVALID_TOKEN)
	}

	if tokenPayload.UserId != userId {
		return c.JSON(http.StatusUnauthorized, UNAUTHORIZED)
	}

	switch data.Field {
	case "username":
		updatedUserParam := db.UpdateUsernameParams{
			Username: data.Payload,
			ID:       userId,
		}

		user, err := s.store.UpdateUsername(c.Request().Context(), updatedUserParam)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, newError(err.Error()))
		}

		return c.JSON(http.StatusOK, newResponse(user))
	default:
		return c.JSON(http.StatusBadRequest, newError("i don't know what you want to update"))
	}
}

func (s *Server) deleteUser(c echo.Context) error {
	return c.JSON(200, "delete User")
}
