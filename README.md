# Echo Extensions

A collection of opinionated extensions and utilities for the [Echo](https://echo.labstack.com/) web framework, providing useful scaffolding and additional functionality to make Echo development easier and more productive.

## Features

- **Swagger Integration**: Easily configure and mount Swagger documentation for your Echo applications
- **Prefix Handling**: Utilities for properly formatting URL prefixes for various extension endpoints
- **Healthcheck Endpoint**: Automatic healthcheck endpoint configuration
- **Custom Middleware**: Pre-configured middlewares for logging, CORS, and recovery
- **Flexible Routing**: Simple group-based routing with middleware support
- **Environment Awareness**: Different behavior based on environment (production vs development)

## Configuration Options

### ServerConfig

Main configuration for the Echo server instance.

| Option | Description | Default Value |
|--------|-------------|---------------|
| PathPrefix | Base path prefix for all routes | `/` |
| Host | Server host address | `0.0.0.0` |
| Port | Server port | `8080` |
| HealthcheckPath | Path suffix for the healthcheck endpoint | `/healthcheck` |
| SkipPaths | Paths to skip for certain middleware (e.g., logging) | `["/", "/healthcheck"]` |
| ExtraCORSHeaders | Additional CORS headers to include beyond the defaults | `[]` |
| SwaggerConfig | Swagger documentation configuration | See below |

### SwaggerConfig

Configuration options for Swagger documentation integration.

| Option | Description | Default Value |
|--------|-------------|---------------|
| Prefix | URL path prefix for accessing Swagger documentation | `/docs` |

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| APP_ENV | Application environment (local, development, production) | `local` |
| SWAGGER_CREDENTIALS | Basic auth credentials for Swagger docs in format `username:password` | `:` |

## Usage

```go
// Basic server with default configuration
server := echoext.New()
server.Start()

// Server with custom configuration
config := echoext.ServerConfig{
    PathPrefix:      "/api/v1",
    Host:            "localhost",
    Port:            3000,
    HealthcheckPath: "/health",
    SkipPaths:       []string{"/metrics", "/status"},
    ExtraCORSHeaders: []string{"X-Api-Key", "X-Custom-Header"},
    SwaggerConfig: echoext.SwaggerConfig{
        Prefix: "/swagger",
    },
}
server := echoext.New(config)

// Creating route groups
server.Group("/users", func(g *echo.Group) {
    g.GET("", listUsers)
    g.POST("", createUser)
    g.GET("/:id", getUser)
})

// Access underlying echo instance
e := server.Engine()

server.Start()
```

## Middleware

The server comes preconfigured with several middleware:

- **Logger**: Logs HTTP requests with customizable path skipping
- **Recovery**: Recovers from panics and returns 500 internal server error
- **CORS**: Configures Cross-Origin Resource Sharing with sensible defaults
  - Default headers include: `Content-Type`, `Content-Length`, `Accept-Encoding`, `X-CSRF-Token`, `Authorization`, `accept`, `origin`, `Cache-Control`, `X-Requested-With`
  - Can be extended with custom headers via the `ExtraCORSHeaders` configuration option


