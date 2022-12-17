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
	// swagger:model
	createPostRequest struct {
		Content        string    `json:"content" validate:"required,max=10000"`
		UserIdentityId uuid.UUID `json:"user_identity_id" validate:"required"`
	}

	// swagger:model
	updatePostRequest struct {
		Content string `json:"content" validate:"required,max=10000"`
	}
)

// list all posts
func (s *Server) listAllPosts(c echo.Context) error {
	/* create a swagger documentation */
	// swagger:operation GET /posts posts listAllPosts
	// ---
	// summary: List all posts
	// description: List all posts
	// parameters:
	// - name: page
	//   in: query
	//   description: page number
	//   required: false
	//   type: integer
	//   format: int32
	// responses:
	//   '200':
	//     description: OK
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/definitions/SuccessResponse"
	//   '400':
	//     description: Bad Request
	//     schema:
	//       "$ref": "#/definitions/BadRequestResponse"
	//   '500':
	//     description: Internal Server Error
	//     schema:
	//       "$ref": "#/definitions/InternalErrorResponse"

	page := c.QueryParam("page")
	if page == "" {
		page = "0"
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

	posts, err := s.store.ListAllPosts(c.Request().Context(), int32(pageInt*10))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newError(err.Error()))
	}

	return c.JSON(http.StatusOK, newResponse(posts))
}

// get a post by id
func (s *Server) getPostById(c echo.Context) error {
	/* create a swagger documentation */
	// swagger:operation GET /posts/{id} posts getPostById
	// ---
	// summary: Get a post by id
	// description: Get a post by id
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
	//       "$ref": "#/definitions/SuccessResponse"
	//   '400':
	//     description: Bad Request
	//     schema:
	//       "$ref": "#/definitions/BadRequestResponse"
	//   '401':
	//     description: Unauthorized
	//     schema:
	//       "$ref": "#/definitions/UnauthorizedResponse"
	//   '404':
	//     description: Not Found
	//     schema:
	//       "$ref": "#/definitions/NotFoundResponse"
	//   '500':
	//     description: Internal Server Error
	//     schema:
	//       "$ref": "#/definitions/InternalErrorResponse"

	id := c.Param("id")
	postId, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

	post, err := s.store.GetPostById(c.Request().Context(), postId)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, NOT_FOUND)
		}
		return c.JSON(http.StatusInternalServerError, newError(err.Error()))
	}

	return c.JSON(http.StatusOK, newResponse(post))
}

// create a new post
func (s *Server) createNewPost(c echo.Context) error {
	/* create a swagger documentation */
	// swagger:operation POST /posts posts createNewPost
	// ---
	// summary: Create a new post
	// description: Create a new post
	// parameters:
	// - name: body
	//   in: body
	//   description: post
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/createPostRequest"
	// security:
	// - key: []
	//
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

	req := new(createPostRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	post, err := s.store.CreatePost(c.Request().Context(), db.CreatePostParams{
		ID:             uuid.New(),
		Content:        req.Content,
		UserIdentityID: req.UserIdentityId,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, NOT_FOUND)
		}
		return c.JSON(http.StatusInternalServerError, newError(err.Error()))
	}

	return c.JSON(http.StatusOK, newResponse(post))
}

// update a post
func (s *Server) updatePost(c echo.Context) error {
	/* create a swagger documentation */
	// swagger:operation PATCH /posts/{id} posts updatePost
	// ---
	// summary: Update a post
	// description: Update a post
	// parameters:
	// - name: id
	//   in: path
	//   description: post id
	//   required: true
	//   type: string
	// - name: body
	//   in: body
	//   description: post
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/updatePostRequest"
	// security:
	// - key: []
	//
	// responses:
	//   '200':
	//     description: OK
	//     schema:
	//       "$ref": "#/definitions/SuccessResponse"
	//   '400':
	//     description: Bad Request
	//     schema:
	//       "$ref": "#/definitions/BadRequestResponse"
	//   '401':
	//     description: Unauthorized
	//     schema:
	//       "$ref": "#/definitions/UnauthorizedResponse"
	//   '500':
	//     description: Internal Server Error
	//     schema:
	//       "$ref": "#/definitions/InternalErrorResponse"

	req := new(updatePostRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

	tokenPayload, ok := c.Get("user").(*token.Payload)
	if !ok {
		return c.JSON(http.StatusUnauthorized, UNAUTHORIZED)
	}

	userIdentity, err := s.store.GetUserIdentityByUserId(c.Request().Context(), tokenPayload.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, NOT_FOUND)
		}
		return c.JSON(http.StatusInternalServerError, INTERNAL_ERROR)
	}

	// get the post
	post, err := s.store.GetPostById(c.Request().Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, NOT_FOUND)
		}
		return c.JSON(http.StatusInternalServerError, INTERNAL_ERROR)
	}

	// check if the user is the owner of the post
	if post.UserIdentityID != userIdentity.ID {
		return c.JSON(http.StatusUnauthorized, UNAUTHORIZED)
	}

	updatedPost, err := s.store.UpdatePost(c.Request().Context(), db.UpdatePostParams{
		ID:      id,
		Content: req.Content,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newError(err.Error()))
	}

	return c.JSON(http.StatusOK, newResponse(updatedPost))
}

// delete a post
func (s *Server) deletePost(c echo.Context) error {
	/* create a swagger documentation */
	// swagger:operation DELETE /posts/{id} posts deletePost
	// ---
	// summary: Delete a post
	// description: Delete a post
	// parameters:
	// - name: id
	//   in: path
	//   description: post id
	//   required: true
	//   type: string
	// security:
	// - key: []
	//
	// responses:
	//   '200':
	//     description: OK
	//     schema:
	//       "$ref": "#/definitions/SuccessResponse"
	//   '400':
	//     description: Bad Request
	//     schema:
	//       "$ref": "#/definitions/BadRequestResponse"
	//   '401':
	//     description: Unauthorized
	//     schema:
	//       "$ref": "#/definitions/UnauthorizedResponse"
	//   '500':
	//     description: Internal Server Error
	//     schema:
	//       "$ref": "#/definitions/InternalErrorResponse"

	id := c.Param("id")
	postId, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

	tokenPayload, ok := c.Get("user").(*token.Payload)
	if !ok {
		return c.JSON(http.StatusUnauthorized, UNAUTHORIZED)
	}

	userIdentity, err := s.store.GetUserIdentityByUserId(c.Request().Context(), tokenPayload.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, NOT_FOUND)
		}
		return c.JSON(http.StatusInternalServerError, INTERNAL_ERROR)
	}

	// get the post
	post, err := s.store.GetPostById(c.Request().Context(), postId)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, NOT_FOUND)
		}
		return c.JSON(http.StatusInternalServerError, INTERNAL_ERROR)
	}

	// check if the post belongs to the user
	if post.UserIdentityID != userIdentity.ID {
		return c.JSON(http.StatusUnauthorized, UNAUTHORIZED)
	}

	// delete the post
	deletedPost, err := s.store.DeletePost(c.Request().Context(), postId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, INTERNAL_ERROR)
	}

	return c.JSON(http.StatusOK, newResponse(deletedPost))
}
