package handler

import (
	db "cnfs/db/sqlc"
	"cnfs/token"
	"cnfs/utils"
	"database/sql"
	"fmt"
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
		SessionId uuid.UUID `json:"session_id" validate:"required"`
		Field     string    `json:"field" validate:"gte=8,required"`
		Payload   string    `json:"payload" validate:"gte=8,required"`
	}

	deleteUserRequest struct {
		SessionId uuid.UUID `json:"session_id" validate:"required"`
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
		return c.JSON(http.StatusInternalServerError, INTERNAL_ERROR)
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
		return c.JSON(http.StatusInternalServerError, INTERNAL_ERROR)
	}

	return c.JSON(200, newResponse(user))
}

func (s *Server) updateUser(c echo.Context) error {
	id := c.Param("id")
	userId, err := uuid.Parse(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

	var data updateUserRequest

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

	_, err = s.store.GetSessionById(c.Request().Context(), data.SessionId)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusUnauthorized, newError("session expired, please login again"))
		}
		return c.JSON(http.StatusInternalServerError, INTERNAL_ERROR)
	}

	switch data.Field {
	case "username":
		updatedUserParam := db.UpdateUsernameParams{
			Username: data.Payload,
			ID:       userId,
		}

		user, err := s.store.UpdateUsername(c.Request().Context(), updatedUserParam)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, INTERNAL_ERROR)
		}

		sessionRemoved, err := s.store.DeleteSession(c.Request().Context(), data.SessionId)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.JSON(http.StatusBadRequest, newError("session doesn't exist"))
			}
			return c.JSON(http.StatusInternalServerError, INTERNAL_ERROR)
		}

		return c.JSON(http.StatusOK, newResponse(fmt.Sprintf("user %s's username has been updated, session %s is deleted, please login again.", user, sessionRemoved)))
	case "password":
		hashedPassword, err := utils.HashPassword(data.Payload)
		if err != nil {
			return c.JSON(http.StatusBadRequest, newError(err.Error()))
		}

		updateUserParam := db.UpdateUserPasswordParams{
			ID:       userId,
			Password: hashedPassword,
		}

		user, err := s.store.UpdateUserPassword(c.Request().Context(), updateUserParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, newError(err.Error()))
		}

		sessionRemoved, err := s.store.DeleteSession(c.Request().Context(), data.SessionId)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.JSON(http.StatusBadRequest, newError("session doesn't exist"))
			}
			return c.JSON(http.StatusInternalServerError, INTERNAL_ERROR)
		}

		return c.JSON(http.StatusOK, newResponse(fmt.Sprintf("user %s's password has been updated, session %s is deleted, please login again.", user, sessionRemoved)))
	default:
		return c.JSON(http.StatusBadRequest, newError("i don't know what you want to update"))
	}
}

func (s *Server) deleteUser(c echo.Context) error {
	id := c.Param("id")
	userId, err := uuid.Parse(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, MISSING_FIELD)
	}

	var data deleteUserRequest

	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, MISSING_FIELD)
	}

	if err := c.Validate(&data); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	_, err = s.store.GetSessionById(c.Request().Context(), data.SessionId)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusUnauthorized, UNAUTHORIZED)
		}

		return c.JSON(http.StatusInternalServerError, INTERNAL_ERROR)
	}

	tokenPayload, ok := c.Get("user").(*token.Payload)
	if !ok {
		return c.JSON(http.StatusUnauthorized, INVALID_TOKEN)
	}

	if tokenPayload.UserId != userId {
		return c.JSON(http.StatusUnauthorized, UNAUTHORIZED)
	}

	user, err := s.store.DeleteOneUser(c.Request().Context(), userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusUnauthorized, UNAUTHORIZED)
		}

		return c.JSON(http.StatusInternalServerError, INTERNAL_ERROR)
	}

	_, err = s.store.DeleteSessionByUserId(c.Request().Context(), userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusUnauthorized, UNAUTHORIZED)
		}

		return c.JSON(http.StatusInternalServerError, INTERNAL_ERROR)
	}

	return c.JSON(http.StatusOK, newResponse(fmt.Sprintf(`user %s has been deleted, all sessions has been revoked`, user)))
}
