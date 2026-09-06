package echoext

import "github.com/labstack/echo/v4"

type HandlerFunc func(c Context) error
type MiddlewareFunc func(next HandlerFunc) HandlerFunc

type Group struct {
	*echo.Group
}

// adaptMiddleware converts our custom middleware to echo middleware
func adaptMiddleware(m MiddlewareFunc) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Unwrap back to the echo context the next handler expects.
			wrappedNext := func(ctx Context) error {
				return next(ctx.Context)
			}

			return m(wrappedNext)(Context{c})
		}
	}
}

func (g *Group) NewGroup(prefix string, middlewares ...MiddlewareFunc) *Group {
	p := escapePath(prefix)

	// Convert our middleware to echo middleware
	echoMiddlewares := make([]echo.MiddlewareFunc, len(middlewares))
	for i, m := range middlewares {
		echoMiddlewares[i] = adaptMiddleware(m)
	}

	return &Group{g.Group.Group(p, echoMiddlewares...)}
}

// applyMiddleware wraps the handler with all middleware functions
func applyMiddleware(h HandlerFunc, middleware ...MiddlewareFunc) HandlerFunc {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}

// adaptHandler converts our custom handler to an echo.HandlerFunc
func adaptHandler(h HandlerFunc, middleware ...MiddlewareFunc) echo.HandlerFunc {
	// Apply all middleware to the handler
	handler := applyMiddleware(h, middleware...)

	// Convert to echo.HandlerFunc
	return func(c echo.Context) error {
		return handler(Context{c})
	}
}

// GET registers a new GET route for the group with a custom Context handler.
func (g *Group) GET(path string, h HandlerFunc, m ...MiddlewareFunc) *echo.Route {
	return g.Group.GET(path, adaptHandler(h, m...))
}

// POST registers a new POST route for the group with a custom Context handler.
func (g *Group) POST(path string, h HandlerFunc, m ...MiddlewareFunc) *echo.Route {
	return g.Group.POST(path, adaptHandler(h, m...))
}

// PUT registers a new PUT route for the group with a custom Context handler.
func (g *Group) PUT(path string, h HandlerFunc, m ...MiddlewareFunc) *echo.Route {
	return g.Group.PUT(path, adaptHandler(h, m...))
}

// DELETE registers a new DELETE route for the group with a custom Context handler.
func (g *Group) DELETE(path string, h HandlerFunc, m ...MiddlewareFunc) *echo.Route {
	return g.Group.DELETE(path, adaptHandler(h, m...))
}

// PATCH registers a new PATCH route for the group with a custom Context handler.
func (g *Group) PATCH(path string, h HandlerFunc, m ...MiddlewareFunc) *echo.Route {
	return g.Group.PATCH(path, adaptHandler(h, m...))
}
