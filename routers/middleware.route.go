package routers

import (
	"os"

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

func AuthMiddleware() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey:    []byte(os.Getenv("SECRET_KEY")),
	})
}
