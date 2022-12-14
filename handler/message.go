package handler

import (
	db "cnfs/db/sqlc"
	"cnfs/token"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type (
	createMessageRequest struct {
		ReceiverId uuid.UUID `json:"receiver_id" validate:"required"`
		Content    string    `json:"content" validate:"required,gte=1"`
	}
)

func (s *Server) createMessage(c echo.Context) error {
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
	}
	message, err := s.store.CreateMessage(c.Request().Context(), msgArg)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

	return c.JSON(200, newResponse(message))
}

func (s *Server) listMessages(c echo.Context) error {
	id := c.Param("id")
	userId, err := uuid.Parse(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
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

func (s *Server) getMessageById(c echo.Context) error {
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

	updateStatusParam := db.UpdateMessageStatusParams{
		ID:         msg.ID,
		ReceiverID: msg.ReceiverID,
	}
	_, err = s.store.UpdateMessageStatus(c.Request().Context(), updateStatusParam)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, NOT_FOUND)
		}
		return c.JSON(http.StatusInternalServerError, INTERNAL_ERROR)
	}

	return c.JSON(200, newResponse(msg))
}

func (s *Server) deleteMessage(c echo.Context) error {
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
