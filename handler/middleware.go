package handler

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"golang.org/x/time/rate"
)

func (s *Server) loggerMiddleware(logger *zerolog.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:     true,
		LogStatus:  true,
		LogLatency: true,
		LogHost:    true,
		LogMethod:  true,
		LogError:   true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			status_code := v.Status
			var log_row *zerolog.Event

			if v.Error != nil {
				log_row = logger.Error().Err(v.Error)
			} else {
				log_row = logger.Info()
			}

			log_row.
				Str("host", v.Host).
				Time("time", v.StartTime.UTC()).
				Str("URI", v.URI).
				Int("status", status_code).
				Str("method", v.Method).
				Str("latency", v.Latency.String()).
				Msg("request")

			return nil
		},
	})
}

func (s *Server) corsMiddleware() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			s.cfg.ClientUrl(),
		},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
		},
		MaxAge: 86400,
	})
}

func (s *Server) rateLimitMiddleware(limit rate.Limit) echo.MiddlewareFunc {
	return middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(limit))
}

const (
	authorizationHeaderKey  = "authorization"
	authorizationHeaderType = "bearer"
	authorizationPayloadKey = "user"
)

func (s *Server) authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authorizationHeader := c.Request().Header.Get(authorizationHeaderKey)
		if authorizationHeader == "" {
			return c.JSON(http.StatusUnauthorized, newError("authorization header is missing"))
		}

		fields := strings.Fields(authorizationHeader)

		if len(fields) < 2 {
			return c.JSON(http.StatusUnauthorized, newError("invalid authorization header format"))
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationHeaderType {
			return c.JSON(http.StatusUnauthorized, newError("unsupported authorization header type"))
		}

		accessToken := fields[1]
		payload, err := s.tokenMaker.VerifyToken(accessToken)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, newError(err.Error()))
		}

		c.Set(authorizationPayloadKey, payload)
		return next(c)
	}
}
