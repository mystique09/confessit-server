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
		AllowOrigins: []string{os.Getenv("FRONTEND_URL")},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	})
}

func CustomCSRFMiddleware() echo.MiddlewareFunc {
	return middleware.CSRFWithConfig(middleware.DefaultCSRFConfig)
}

func CustomRateLimitMiddleware() echo.MiddlewareFunc {
	return middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20))
}
