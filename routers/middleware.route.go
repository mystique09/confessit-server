package routers

import (
	"confessit/models"
	"confessit/utils"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CustomLoggerMiddleware() echo.MiddlewareFunc {
	return middleware.LoggerWithConfig(
		middleware.LoggerConfig{
			Format: `[${time_rfc3339}] ${status} /${method} ${host}${path} ${latency_human}` + "\n",
		},
	)
}

func CustomCORSMiddleware() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{os.Getenv("FRONTEND_URL"), "confessit.vercel.app", "https://confessit.vercel.app"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
	})
}

func CustomCSRFMiddleware() echo.MiddlewareFunc {
	return middleware.CSRFWithConfig(middleware.DefaultCSRFConfig)
}

func CustomRateLimitMiddleware() echo.MiddlewareFunc {
	return middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20))
}

func (r *Route) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := c.Cookie("session_token")

		if err != nil {
			if err == http.ErrNoCookie {
				return c.JSON(http.StatusUnauthorized, "no cookie in headers")
			}
			return c.JSON(http.StatusBadRequest, "err getting cookie in headers")
		}

		sessionToken := session.Value
		var token models.Session

		r.Conn.Model(&models.Session{}).Where("id = ?", sessionToken).Find(&token)

		if token.ID == uuid.Nil {
			return c.JSON(http.StatusUnauthorized, "Unauthorized")
		}

		if token.IsExpired() {
			new_cookie := utils.CreateCookie("session_token", "", -1)
			c.SetCookie(&new_cookie)
			r.Conn.Delete(&models.Session{}, "id = ?", sessionToken)

			return c.JSON(http.StatusUnauthorized, "expired cookie")
		}

		return next(c)
	}
}
