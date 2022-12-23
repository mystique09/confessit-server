package handler

import (
	"cnfs/common"
	db "cnfs/db/sqlc"
	"database/sql"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type (
	// swagger:model
	loginRequest struct {
		// The username
		// required: true
		Username string `json:"username" validate:"required,gte=8"`
		// The password
		// required: true
		Password string `json:"password" validate:"required,gte=8"`
	}

	// swagger:model
	logoutRequest struct {
		SessionId uuid.UUID `json:"session_id" validate:"required"`
	}

	// swagger:model loginResponse
	loginResponse struct {
		// The session id that is saved in the db
		SessionId uuid.UUID `json:"session_id"`
		// The access token that is used to access resources
		AccessToken string `json:"access_token"`
		// The expiration date of access token
		AccessTokenExpiresAt time.Time `json:"access_token_expiry"`
		// The refresh token that is used to refresh access token
		RefreshToken string `json:"refresh_token"`
		// The expiration date of refresh token
		RefreshTokenExpiresAt time.Time `json:"refresh_token_expiry"`
		// The user information needed for client
		User         db.User         `json:"user"`
		UserIdentity db.UserIdentity `json:"user_identity"`
	}
)

// Login endpoint.
func (s *Server) loginUser(c echo.Context) error {
	// The login handler.
	// swagger:operation POST /auth auth loginUser
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
	//   description: payload needed for login
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/loginRequest"
	//
	// responses:
	//  200:
	//	  description: Login successful
	//	  schema:
	//	     type: object
	//		 	"$ref": "#/definitions/loginResponse"
	//  400:
	//	  description: Bad request
	//	  schema:
	//	     type: object
	//		 	"$ref": "#/definitions/BadRequestResponse"
	//  500:
	//	  description: Internal server error
	//	  schema:
	//	     type: object
	//		 	"$ref": "#/definitions/InternalErrorResponse"

	var data loginRequest

	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

	if err := c.Validate(&data); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	user, err := s.store.GetUserByUsername(c.Request().Context(), data.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(404, newError("user not found"))
		}
		return c.JSON(500, newError(err.Error()))
	}

	if err := common.CheckPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		return c.JSON(403, newError("password mismatch"))
	}

	accessToken, accessTokenPayload, err := s.tokenMaker.CreateToken(user.ID, user.Username, s.cfg.AccessTokenDuration)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

	refreshToken, refreshTokenPayload, err := s.tokenMaker.CreateToken(user.ID, user.Username, s.cfg.RefreshTokenDuration)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

	// set refresh token as cookie
	cookie := new(http.Cookie)
	cookie.Path = "/"
	cookie.Domain = c.Request().URL.String()
	cookie.Name = "refresh_token_cookie"
	cookie.Value = refreshToken
	cookie.MaxAge = refreshTokenPayload.ExpiredAt.Second()
	cookie.HttpOnly = true
	cookie.Secure = true
	c.SetCookie(cookie)

	newSessionArg := db.CreateSessionParams{
		ID:           refreshTokenPayload.Id,
		UserID:       user.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    c.Request().UserAgent(),
		ClientIp:     c.RealIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshTokenPayload.ExpiredAt,
	}

	newSession, err := s.store.CreateSession(c.Request().Context(), newSessionArg)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusBadRequest, newError(err.Error()))
		}
		return c.JSON(500, newError(err.Error()))
	}

	userIdentity, err := s.store.GetUserIdentityByUserId(c.Request().Context(), user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusBadRequest, newError(err.Error()))
		}
		return c.JSON(500, newError(err.Error()))
	}

	resp := loginResponse{
		SessionId:             newSession.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenPayload.ExpiredAt,
		User:                  user,
		UserIdentity:          userIdentity,
	}

	return c.JSON(200, newResponse(resp))
}

// clears the session in the database
func (s *Server) logoutUser(c echo.Context) error {
	/* create a swagger documentation*/
	// logout user
	//
	// swagger:operation DELETE /auth/clear auth logoutUser
	//
	// ---
	//
	// consumes:
	// - application/json
	//
	// produces:
	// - application/json
	//
	// parameters:
	// - name: body
	//   in: body
	//   description: payload needed for logout
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/logoutRequest"
	//
	// responses:
	//  200:
	//	  description: Logout successful
	//	  schema:
	//	     type: object
	//		 	"$ref": "#/definitions/SuccessResponse"
	//  400:
	//	  description: Bad request
	//	  schema:
	//	     type: object
	//		 	"$ref": "#/definitions/BadRequestResponse"
	//  500:
	//	  description: Internal server error
	//	  schema:
	//	     type: object
	//		 	"$ref": "#/definitions/InternalErrorResponse"

	var data logoutRequest

	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

	if err := c.Validate(&data); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	deletedSession, err := s.store.DeleteSession(c.Request().Context(), data.SessionId)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, NOT_FOUND)
		}
		return c.JSON(http.StatusInternalServerError, INTERNAL_ERROR)
	}

	return c.JSON(http.StatusOK, newResponse(deletedSession))
}
