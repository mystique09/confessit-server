package handler

import (
	db "cnfs/db/sqlc"
	"cnfs/token"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type (
	// swagger:model
	createMessageRequest struct {
		// the receiver_id
		// type: uuid
		// required: true
		ReceiverId uuid.UUID `json:"receiver_id" validate:"required"`

		// the actual message
		// type: string
		// required: true
		Content string `json:"content" validate:"required,gte=1"`
	}
)

// creates a new message
func (s *Server) createMessage(c echo.Context) error {
	// Creates a new message.
	// swagger:operation POST /messages messages createMessage
	//
	// ---
	// produces:
	// - application/json
	// consumes:
	// - application/json
	//
	// parameters:
	// - name: body
	//   in: body
	//   description: this payload holds the receiver_id and the message
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/createMessageRequest"
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

	data := new(createMessageRequest)

	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, MISSING_FIELD)
	}

	if err := c.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	user, err := s.store.GetUserById(c.Request().Context(), data.ReceiverId)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, NOT_FOUND)
		}
		return c.JSON(http.StatusInternalServerError, INTERNAL_ERROR)
	}

	msgArg := db.CreateMessageParams{
		ID:         uuid.New(),
		ReceiverID: user.ID,
		Content:    data.Content,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	message, err := s.store.CreateMessage(c.Request().Context(), msgArg)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

	return c.JSON(200, newResponse(message))
}

// List all messages that has sent to the user
func (s *Server) listMessages(c echo.Context) error {
	// List all messages that has sent to the user.
	// swagger:operation GET /users/{id}/messages messages listMessages
	//
	// ---
	// produces:
	// - application/json
	// consumes:
	// - application/json
	//
	// parameters:
	// - name: page
	//   in: query
	//   description: the page number
	//   required: false
	//   type: integer
	//   format: int64
	// - name: id
	//   in: path
	//   description: the user id
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

	tokenPayload, ok := c.Get("user").(*token.Payload)
	if !ok {
		return c.JSON(http.StatusUnauthorized, UNAUTHORIZED)
	}

	if tokenPayload.UserId != userId {
		return c.JSON(http.StatusUnauthorized, UNAUTHORIZED)
	}

	pageParam := c.QueryParam("page")
	if pageParam == "" {
		pageParam = "0"
	}

	page, err := strconv.ParseUint(pageParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

	msgParam := db.ListMessageParams{
		ReceiverID: userId,
		Offset:     int32(page * 10),
	}
	messages, err := s.store.ListMessage(c.Request().Context(), msgParam)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, NOT_FOUND)
		}
		return c.JSON(http.StatusInternalServerError, INTERNAL_ERROR)
	}

	return c.JSON(200, newResponse(messages))
}

// Get a message by id
func (s *Server) getMessageById(c echo.Context) error {
	// Get a message by id.
	// swagger:operation GET /messages/{id} messages getMessageById
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
	//   description: the message id
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

	id := c.Param("id")
	msgId, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

	msg, err := s.store.GetMessageById(c.Request().Context(), msgId)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, NOT_FOUND)
		}
		return c.JSON(http.StatusInternalServerError, INTERNAL_ERROR)
	}

	return c.JSON(200, newResponse(msg))
}

func (s *Server) updateMessage(c echo.Context) error {
	// Update a message.
	// swagger:operation PUT /messages/{id} messages updateMessage
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
	//   description: the message id
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

	id := c.Param("id")
	messageId, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

	tokenPayload, ok := c.Get("user").(*token.Payload)
	if !ok {
		return c.JSON(http.StatusUnauthorized, UNAUTHORIZED)
	}

	msg, err := s.store.GetMessageById(c.Request().Context(), messageId)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, NOT_FOUND)
		}
		return c.JSON(http.StatusInternalServerError, INTERNAL_ERROR)
	}

	if msg.ReceiverID != tokenPayload.UserId {
		return c.JSON(http.StatusUnauthorized, UNAUTHORIZED)
	}

	message, err := s.store.UpdateMessageStatus(c.Request().Context(), db.UpdateMessageStatusParams{
		ID:         messageId,
		ReceiverID: tokenPayload.UserId,
		UpdatedAt:  time.Now(),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, NOT_FOUND)
		}
		return c.JSON(http.StatusInternalServerError, INTERNAL_ERROR)
	}

	return c.JSON(200, newResponse(message))
}

func (s *Server) deleteMessage(c echo.Context) error {
	// Delete a message.
	// swagger:operation DELETE /messages/{id} messages deleteMessage
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
	//   description: the message id
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

	id := c.Param("id")
	msgId, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

	tokenPayload, ok := c.Get("user").(*token.Payload)
	if !ok {
		return c.JSON(http.StatusUnauthorized, UNAUTHORIZED)
	}

	msg, err := s.store.GetMessageById(c.Request().Context(), msgId)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, NOT_FOUND)
		}
		return c.JSON(http.StatusInternalServerError, INTERNAL_ERROR)
	}

	if msg.ReceiverID != tokenPayload.UserId {
		return c.JSON(http.StatusUnauthorized, UNAUTHORIZED)
	}

	deleteMsgParam := db.DeleteOneMessageParams{
		ReceiverID: tokenPayload.UserId,
		ID:         msgId,
	}
	deletedMsg, err := s.store.DeleteOneMessage(c.Request().Context(), deleteMsgParam)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, NOT_FOUND)
		}
		return c.JSON(http.StatusInternalServerError, INTERNAL_ERROR)
	}

	return c.JSON(200, newResponse(deletedMsg))
}
