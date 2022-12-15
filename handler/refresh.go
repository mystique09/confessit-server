package handler

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type (
	// swagger:model
	refreshRequest struct {
		// The refresh token
		// required: true
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	// swagger:model
	accessTokenRequest struct {
		// The access token
		// required: true
		AccessToken string `json:"access_token" validate:"required"`
	}

	// swagger:model
	accessTokenResponse struct {
		// The access token
		// in: body
		AccessToken string `json:"access_token"`
		// The access token expiry
		// in: body
		AccessTokenExpiresAt time.Time `json:"access_token_expiry"`
	}
)

func (s *Server) refreshAccessToken(c echo.Context) error {
	// Refresh access token using refresh token.
	//
	// swagger:operation POST /auth/refresh auth refreshAccessTokenRequestBody
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
	//   description: the payload needed for refreshing access token
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/refreshRequest"
	//
	// responses:
	//  200:
	//	  description: Refresh access token successful
	//	  schema:
	//	     type: object
	//		 	"$ref": "#/definitions/accessTokenResponse"
	//  400:
	//	  description: Bad request
	//	  schema:
	//	     type: object
	//		 	"$ref": "#/definitions/BadRequestResponse"
	//  401:
	//	  description: Unauthorized
	//	  schema:
	//	     type: object
	//		 	"$ref": "#/definitions/UnauthorizedResponse"
	//  404:
	//	  description: Not found
	//	  schema:
	//	     type: object
	//		 	"$ref": "#/definitions/NotFoundResponse"
	//  500:
	//	  description: Internal server error
	//	  schema:
	//	     type: object
	//		 	"$ref": "#/definitions/InternalErrorResponse"

	var data refreshRequest

	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

	if err := c.Validate(&data); err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

	refreshTokenPayload, err := s.tokenMaker.VerifyToken(data.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, newError(err.Error()))
	}

	session, err := s.store.GetSessionById(c.Request().Context(), refreshTokenPayload.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, newError(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, newError(err.Error()))
	}

	if session.IsBlocked {
		return c.JSON(http.StatusUnauthorized, newError("blocked session"))
	}

	if (session.UserID != refreshTokenPayload.UserId) || (session.Username != refreshTokenPayload.Username) {
		return c.JSON(http.StatusUnauthorized, newError("invalid session token"))
	}

	if session.RefreshToken != data.RefreshToken {
		return c.JSON(http.StatusUnauthorized, newError("mismatched session token"))
	}

	if time.Now().After(session.ExpiresAt) {
		return c.JSON(http.StatusUnauthorized, newError("expired session"))
	}

	newAccessToken, newAccessTokenPayload, err := s.tokenMaker.CreateToken(session.UserID, session.Username, s.cfg.AccessTokenDuration)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newError(err.Error()))
	}

	resp := accessTokenResponse{
		AccessToken:          newAccessToken,
		AccessTokenExpiresAt: newAccessTokenPayload.ExpiredAt,
	}

	return c.JSON(http.StatusOK, newResponse(resp))
}

func (s *Server) validateAccessToken(c echo.Context) error {
	// Validate access token.
	//
	// swagger:operation POST /auth/validate auth validateAccessTokenRequestBody
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
	//   description: the access token to validate
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/accessTokenRequest"
	//
	// responses:
	//  200:
	//	  description: Validate access token successful
	//	  schema:
	//	     type: object
	//		 	"$ref": "#/definitions/accessTokenResponse"
	//  400:
	//	  description: Bad request
	//	  schema:
	//	     type: object
	//		 	"$ref": "#/definitions/BadRequestResponse"
	//  401:
	//	  description: Unauthorized
	//	  schema:
	//	     type: object
	//		 	"$ref": "#/definitions/UnauthorizedResponse"
	//  500:
	//	  description: Internal server error
	//	  schema:
	//	     type: object
	//		 	"$ref": "#/definitions/InternalErrorResponse"

	var data accessTokenRequest

	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

	if err := c.Validate(&data); err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

	payload, err := s.tokenMaker.VerifyToken(data.AccessToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, newError(err.Error()))
	}

	resp := accessTokenResponse{
		AccessToken:          data.AccessToken,
		AccessTokenExpiresAt: payload.ExpiredAt,
	}

	return c.JSON(200, newResponse(resp))
}
