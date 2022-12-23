package handler

import (
	db "cnfs/db/sqlc"
	"cnfs/token"
	"database/sql"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type (
	// swagger:model
	createCommentRequest struct {
		PostId         uuid.UUID `json:"post_id" validate:"required"`
		Content        string    `json:"content" validate:"required"`
		UserIdentityId uuid.UUID `json:"user_identity_id" validate:"required"`
		ParentId       uuid.UUID `json:"parent_id" validate:"required"`
	}

	// swagger:model
	updateCommentRequest struct {
		Content string `json:"content" validate:"required"`
	}
)

// list all comments of a post
func (s *Server) listAllComments(c echo.Context) error {
	// make a swagger docs

	// swagger:operation GET /posts/{id}/comments comments listAllComments
	// ---
	// summary: List all comments of a post
	// description: List all comments of a post
	// parameters:
	// - name: id
	//   in: path
	//   description: post id
	//   required: true
	//   type: string
	// responses:
	//   '200':
	//     description: OK
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/definitions/SuccessResponse"
	//   '404':
	//     description: Not Found
	//     schema:
	//       "$ref": "#/definitions/BadRequestResponse"
	//   '500':
	//     description: Internal Server Error
	//     schema:
	//       "$ref": "#/definitions/InternalErrorResponse"

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

// get comment by id
func (s *Server) getCommentById(c echo.Context) error {
	// swagger:operation GET /comments/{id} comments getCommentById
	// ---
	// summary: Get comment by id
	// description: Get comment by id
	// parameters:
	// - name: id
	//   in: path
	//   description: comment id
	//   required: true
	//   type: string
	// responses:
	//   '200':
	//     description: OK
	//     schema:
	//       "$ref": "#/definitions/SuccessResponse"
	//   '404':
	//     description: Not Found
	//     schema:
	//       "$ref": "#/definitions/BadRequestResponse"
	//   '500':
	//     description: Internal Server Error
	//     schema:
	//       "$ref": "#/definitions/InternalErrorResponse"

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

// create a comment
func (s *Server) createComment(c echo.Context) error {
	// swagger:operation POST /comments comments createComment
	// ---
	// summary: Create a comment
	// description: Create a comment
	// parameters:
	// - name: body
	//   in: body
	//   description: comment
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/createCommentRequest"
	// security:
	// - key: []
	// responses:
	//   '200':
	//     description: OK
	//     schema:
	//       "$ref": "#/definitions/SuccessResponse"
	//   '400':
	//     description: Bad Request
	//     schema:
	//       "$ref": "#/definitions/BadRequestResponse"
	//   '500':
	//     description: Internal Server Error
	//     schema:
	//       "$ref": "#/definitions/InternalErrorResponse"

	req := new(createCommentRequest)
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

	commentId := uuid.New()
	if req.ParentId == uuid.Nil {
		req.ParentId = commentId
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	comment, err := s.store.CreateComment(c.Request().Context(), db.CreateCommentParams{
		ID:             commentId,
		PostID:         req.PostId,
		Content:        req.Content,
		UserIdentityID: req.UserIdentityId,
		ParentID:       req.ParentId,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})

	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, NOT_FOUND)
		}
		return c.JSON(http.StatusInternalServerError, INTERNAL_ERROR)
	}

	return c.JSON(http.StatusOK, newResponse(comment))
}

func (s *Server) updateComment(c echo.Context) error {
	// swagger:operation PUT /comments/{id} comments updateComment
	// ---
	// summary: Update a comment
	// description: Update a comment
	// parameters:
	// - name: id
	//   in: path
	//   description: comment id
	//   required: true
	//   type: string
	// - name: body
	//   in: body
	//   description: comment
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/updateCommentRequest"
	// security:
	// - key: []
	// responses:
	//   '200':
	//     description: OK
	//     schema:
	//       "$ref": "#/definitions/SuccessResponse"
	//   '400':
	//     description: Bad Request
	//     schema:
	//       "$ref": "#/definitions/BadRequestResponse"
	//   '500':
	//     description: Internal Server Error
	//     schema:
	//       "$ref": "#/definitions/InternalErrorResponse"

	req := new(updateCommentRequest)
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

	id := c.Param("id")
	commentId, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(400, newError(err.Error()))
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	tokenPayload, ok := c.Get("user").(*token.Payload)
	if !ok {
		return c.JSON(http.StatusUnauthorized, UNAUTHORIZED)
	}

	userIdentityId, err := s.store.GetUserIdentityByUserId(c.Request().Context(), tokenPayload.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, NOT_FOUND)
		}
		return c.JSON(http.StatusInternalServerError, INTERNAL_ERROR)
	}

	hasComment, err := s.store.GetComment(c.Request().Context(), commentId)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, NOT_FOUND)
		}
		return c.JSON(http.StatusInternalServerError, INTERNAL_ERROR)
	}

	if userIdentityId.ID != hasComment.UserIdentityID {
		return c.JSON(http.StatusUnauthorized, UNAUTHORIZED)
	}

	comment, err := s.store.UpdateComment(c.Request().Context(), db.UpdateCommentParams{
		ID:        commentId,
		Content:   req.Content,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, NOT_FOUND)
		}
		return c.JSON(http.StatusInternalServerError, INTERNAL_ERROR)
	}

	return c.JSON(http.StatusOK, newResponse(comment))
}

func (s *Server) deleteComment(c echo.Context) error {
	// swagger:operation DELETE /comments/{id} comments deleteComment
	// ---
	// summary: Delete a comment
	// description: Delete a comment
	// parameters:
	// - name: id
	//   in: path
	//   description: comment id
	//   required: true
	//   type: string
	// security:
	// - key: []
	// responses:
	//   '200':
	//     description: OK
	//     schema:
	//       "$ref": "#/definitions/SuccessResponse"
	//   '400':
	//     description: Bad Request
	//     schema:
	//       "$ref": "#/definitions/BadRequestResponse"
	//   '500':
	//     description: Internal Server Error
	//     schema:
	//       "$ref": "#/definitions/InternalErrorResponse"

	id := c.Param("id")
	commentId, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(400, newError(err.Error()))
	}

	tokenPayload, ok := c.Get("user").(*token.Payload)
	if !ok {
		return c.JSON(http.StatusUnauthorized, UNAUTHORIZED)
	}

	userIdentityId, err := s.store.GetUserIdentityByUserId(c.Request().Context(), tokenPayload.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, NOT_FOUND)
		}
		return c.JSON(http.StatusInternalServerError, INTERNAL_ERROR)
	}

	hasComment, err := s.store.GetComment(c.Request().Context(), commentId)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, NOT_FOUND)
		}
		return c.JSON(http.StatusInternalServerError, INTERNAL_ERROR)
	}

	if userIdentityId.ID != hasComment.UserIdentityID {
		return c.JSON(http.StatusUnauthorized, UNAUTHORIZED)
	}

	comment, err := s.store.DeleteComment(c.Request().Context(), commentId)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, NOT_FOUND)
		}
		return c.JSON(http.StatusInternalServerError, INTERNAL_ERROR)
	}

	return c.JSON(http.StatusOK, newResponse(comment))
}
