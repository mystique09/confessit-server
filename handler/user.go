package handler

import (
	"cnfs/common"
	db "cnfs/db/sqlc"
	"cnfs/token"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type (
	// swagger:model
	// type: uuid
	// The id of a user in uuid format
	UserId = uuid.UUID

	// swagger:model
	createUserRequest struct {
		// the username
		//
		// unique: true
		// required: true
		Username string `json:"username" validate:"gte=8,required"`
		// the password
		//
		// required: true
		Password string `json:"password" validate:"gte=8,required"`
	}

	// swagger:model
	updateUserRequest struct {
		// SessionId of the user in payload
		//
		// unique: true
		// in: body
		// type: uuid
		SessionId uuid.UUID `json:"session_id" validate:"required"`

		// Field of the user in payload
		//
		// unique: true
		// in: body
		// type: string
		Field string `json:"field" validate:"gte=8,required"`

		// Payload of the user in payload
		//
		// unique: true
		// in: body
		// type: string
		Payload string `json:"payload" validate:"gte=8,required"`
	}

	// swagger:model
	deleteUserRequest struct {
		// SessionId of the user in payload
		//
		// unique: true
		// in: body
		// type: uuid
		SessionId uuid.UUID `json:"session_id" validate:"required"`
	}
)

// Creatas a new user.
func (s *Server) createUser(c echo.Context) error {
	// Create a new user and return the user information.
	// swagger:operation POST /users users createUser
	//
	// ---
	// consumes:
	// - application/json
	//
	// produces:
	// - application/json
	//
	// parameters:
	// - name: body
	//   in: body
	//   description: user information
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/createUserRequest"
	//
	// responses:
	//  200:
	//	  description: Success response with user information.
	//	  schema:
	//	     type: object
	//		 	"$ref": "#/definitions/SuccessResponse"
	//  400:
	//	  description: Bad request.
	//	  schema:
	//	     type: object
	//		 	"$ref": "#/definitions/BadRequestResponse"
	//  500:
	//	  description: Internal error.
	//	  schema:
	//	     type: object
	//		 	"$ref": "#/definitions/InternalErrorResponse"

	var data createUserRequest

	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

	if err := c.Validate(&data); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	hashedPassword, err := common.HashPassword(data.Password)
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

// List all users.
func (s *Server) listUsers(c echo.Context) error {
	// List all users in the system.
	// swagger:operation GET /users users listUsers
	//
	// ---
	//
	// produces:
	// - application/json
	//
	// parameters:
	// - name: page
	//   in: query
	//   description: page number
	//   required: false
	//   type: integer
	//   format: int64
	//
	// security:
	// - key: []
	//
	// responses:
	//  200:
	//	  description: Success response with user information.
	//	  schema:
	//	     type: array
	//		 	items:
	//		 		"$ref": "#/definitions/User"
	//  400:
	//	  description: Bad request.
	//	  schema:
	//	     type: object
	//		 	"$ref": "#/definitions/BadRequestResponse"
	//  500:
	//	  description: Internal error.
	//	  schema:
	//	     type: object
	//		 	"$ref": "#/definitions/InternalErrorResponse"
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

// Get user by id.
func (s *Server) getUserById(c echo.Context) error {
	// Get user by id.
	// swagger:operation GET /users/{id} users getUserById
	//
	// ---
	//
	// produces:
	// - application/json
	//
	// parameters:
	// - name: id
	//   in: path
	//   description: user id
	//   required: true
	//   type: string
	//   format: uuid
	//
	// security:
	// - key: []
	//
	// responses:
	//  200:
	//	  description: Success response with user information.
	//	  schema:
	//	     type: object
	//		 	"$ref": "#/definitions/User"
	//  400:
	//	  description: Bad request.
	//	  schema:
	//	     type: object
	//		 	"$ref": "#/definitions/BadRequestResponse"
	//  500:
	//	  description: Internal error.
	//	  schema:
	//	     type: object
	//		 	"$ref": "#/definitions/InternalErrorResponse"
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

// Update user by id.
func (s *Server) updateUser(c echo.Context) error {
	// Update user by id.
	// swagger:operation PATCH /users/{id} users updateUserById
	//
	// ---
	// consumes:
	// - application/json
	//
	// produces:
	// - application/json
	//
	// parameters:
	// - name: id
	//   in: path
	//   description: user id
	//   required: true
	//   type: string
	//   format: uuid
	// - name: body
	//   in: body
	//   description: user information
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/updateUserRequest"
	//
	// security:
	// - key: []
	//
	// responses:
	//  200:
	//	  description: A success response that says "user's username/password has been updated"
	//	  schema:
	//	     type: object
	//		 	"$ref": "#/definitions/SuccessResponse"
	//  400:
	//	  description: Bad request.
	//	  schema:
	//	     type: object
	//		 	"$ref": "#/definitions/BadRequestResponse"
	//  401:
	//	  description: Unauthorized.
	//	  schema:
	//	     type: object
	//		 	"$ref": "#/definitions/UnauthorizedResponse"
	//  500:
	//	  description: Internal error.
	//	  schema:
	//	     type: object
	//		 	"$ref": "#/definitions/InternalErrorResponse"

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
		hashedPassword, err := common.HashPassword(data.Payload)
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

// Delete user
func (s *Server) deleteUser(c echo.Context) error {
	// Delete one user by id
	// swagger:operation DELETE /users/{id} users deleteUser
	//
	// ---
	// produces:
	// - application/json
	// consumes:
	// - application/json
	//
	// parameters:
	// - name: id
	//   in: path
	//   description: user id
	//   required: true
	//   type: string
	//   format: uuid
	// - name: body
	//   in: body
	//   description: session id
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/deleteUserRequest"
	//
	// security:
	// - key: []
	//
	// responses:
	//  200:
	//	  description: User deleted.
	//	  schema:
	//	     type: object
	//		 	"$ref": "#/definitions/SuccessResponse"
	//  400:
	//	  description: Bad request.
	//	  schema:
	//	     type: object
	//		 	"$ref": "#/definitions/BadRequestResponse"
	//  401:
	//	  description: Unauthorized.
	//	  schema:
	//	     type: object
	//		 	"$ref": "#/definitions/UnauthorizedResponse"
	//  500:
	//	  description: Internal server error.
	//	  schema:
	//	     type: object
	//		 	"$ref": "#/definitions/InternalErrorResponse"
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
		log.Println(err.Error())
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
		log.Println(tokenPayload.UserId, userId)
		return c.JSON(http.StatusUnauthorized, UNAUTHORIZED)
	}

	user, err := s.store.DeleteOneUser(c.Request().Context(), userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusUnauthorized, UNAUTHORIZED)
		}

		return c.JSON(http.StatusInternalServerError, INTERNAL_ERROR)
	}

	return c.JSON(http.StatusOK, newResponse(fmt.Sprintf(`user %s has been deleted, all sessions has been revoked`, user)))
}
