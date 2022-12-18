package handler

import (
	db "cnfs/db/sqlc"
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type (
	createCommentRequest struct {
		PostId         uuid.UUID `json:"post_id" validate:"required"`
		Content        string    `json:"content" validate:"required"`
		UserIdentityId uuid.UUID `json:"user_identity_id" validate:"required"`
		ParentId       uuid.UUID `json:"parent_id"`
	}
)

func (s *Server) listAllComments(c echo.Context) error {
	id := c.Param("id")
	postId, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(400, newError(err.Error()))
	}

	comments, err := s.store.ListAllComments(c.Request().Context(), postId)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, NOT_FOUND)
		}
		return c.JSON(http.StatusInternalServerError, INTERNAL_ERROR)
	}

	return c.JSON(http.StatusOK, newResponse(comments))
}

func (s *Server) getCommentById(c echo.Context) error {
	id := c.Param("id")
	commentId, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(400, newError(err.Error()))
	}

	comment, err := s.store.GetComment(c.Request().Context(), commentId)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, NOT_FOUND)
		}
		return c.JSON(http.StatusInternalServerError, INTERNAL_ERROR)
	}

	return c.JSON(http.StatusOK, newResponse(comment))
}

func (s *Server) createComment(c echo.Context) error {
	var req createCommentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	comment, err := s.store.CreateComment(c.Request().Context(), db.CreateCommentParams{
		ID:             uuid.New(),
		PostID:         req.PostId,
		Content:        req.Content,
		UserIdentityID: req.UserIdentityId,
		ParentID: uuid.NullUUID{
			UUID: req.ParentId,
		},
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, NOT_FOUND)
		}
		return c.JSON(http.StatusInternalServerError, INTERNAL_ERROR)
	}

	return c.JSON(http.StatusOK, newResponse(comment))
}
