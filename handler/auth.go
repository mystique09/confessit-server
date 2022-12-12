package handler

import (
	db "cnfs/db/sqlc"
	"cnfs/utils"
	"database/sql"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type (
	loginRequest struct {
		Username string `json:"username" validate:"required,gte=8"`
		Password string `json:"password" validate:"required,gte=8"`
	}

	loginResponse struct {
		SessionId             uuid.UUID `json:"session_id"`
		AccessToken           string    `json:"access_token"`
		AccessTokenExpiresAt  time.Time `json:"access_token_expiry"`
		RefreshToken          string    `json:"refresh_token"`
		RefreshTokenExpiresAt time.Time `json:"refresh_token_expiry"`
		User                  db.User   `json:"user"`
	}
)

func (s *Server) loginUser(c echo.Context) error {
	var data loginRequest

	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

	if err := c.Validate(&data); err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

	user, err := s.store.GetUserByUsername(c.Request().Context(), data.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(404, newError("user not found"))
		}
		return c.JSON(500, newError(err.Error()))
	}

	if err := utils.CheckPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		return c.JSON(403, newError("password mismatch"))
	}

	accessToken, accessTokenPayload, err := s.tokenMaker.CreateToken(user.ID, user.Username, s.cfg.AccessTokenDuration)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

	refreshToken, refreshTokenPayload, err := s.tokenMaker.CreateToken(user.ID, user.Username, s.cfg.AccessTokenDuration)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newError(err.Error()))
	}

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

	resp := loginResponse{
		SessionId:             newSession.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenPayload.ExpiredAt,
		User:                  user,
	}

	return c.JSON(200, newResponse(resp))
}
