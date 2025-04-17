package echoext

import (
	"net/http"
	"slices"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var CustomRecovery = middleware.RecoverWithConfig(middleware.DefaultRecoverConfig)

var CustomCORS = middleware.CORSWithConfig(middleware.CORSConfig{
	AllowOrigins:     []string{"*"},
	AllowCredentials: true,
	AllowMethods: []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
	},
	AllowHeaders: []string{
		"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token",
		"Authorization", "accept", "origin", "Cache-Control",
		"X-Requested-With",
	},
})

func CustomLogger(c ServerConfig) echo.MiddlewareFunc {
	return middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(ctx echo.Context) bool {
			if ctx.Request().Method == http.MethodOptions {
				return true
			}

			paths := c.escapeSkipPaths()
			if strings.HasPrefix(ctx.Request().URL.Path, c.swaggerPath()) {
				return true
			}

			return slices.Contains(paths, strings.ToLower(ctx.Request().URL.Path))
		},
	})
}
