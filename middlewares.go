package echoext

import (
	"net/http"
	"slices"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var CustomRecovery = middleware.RecoverWithConfig(middleware.DefaultRecoverConfig)

// Default CORS headers used by the middleware
var defaultCORSHeaders = []string{
	"Content-Type", "Content-Length", "Accept-Encoding",
	"X-CSRF-Token", "Authorization", "accept", "origin",
	"Cache-Control", "X-Requested-With",
}

// CustomCORS creates a CORS middleware with the provided config
func CustomCORS(c ServerConfig) echo.MiddlewareFunc {
	// Combine default headers with any extra headers from config
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodPatch,
		},
		AllowHeaders: append(defaultCORSHeaders, c.ExtraCORSHeaders...),
	})
}

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
