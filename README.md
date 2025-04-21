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

## Extended Context

The extension provides an enhanced Context interface that extends Echo's standard Context with additional type-safe getter methods. These methods simplify the retrieval of typed values from context storage.

### Type-Safe Getter Methods

| Method | Return Type | Description |
|--------|-------------|-------------|
| `GetString(key string)` | `string` | Retrieves a string value from context storage |
| `GetBool(key string)` | `bool` | Retrieves a boolean value from context storage |
| `GetInt(key string)` | `int` | Retrieves an int value from context storage |
| `GetInt8(key string)` | `int8` | Retrieves an int8 value from context storage |
| `GetInt16(key string)` | `int16` | Retrieves an int16 value from context storage |
| `GetInt32(key string)` | `int32` | Retrieves an int32 value from context storage |
| `GetInt64(key string)` | `int64` | Retrieves an int64 value from context storage |
| `GetUint(key string)` | `uint` | Retrieves a uint value from context storage |
| `GetUint8(key string)` | `uint8` | Retrieves a uint8 value from context storage |
| `GetUint16(key string)` | `uint16` | Retrieves a uint16 value from context storage |
| `GetUint32(key string)` | `uint32` | Retrieves a uint32 value from context storage |
| `GetUint64(key string)` | `uint64` | Retrieves a uint64 value from context storage |
| `GetFloat64(key string)` | `float64` | Retrieves a float64 value from context storage |

Each method automatically performs type assertion on the value stored in context, returning the zero value of the respective type if the value is not of the expected type or not found.

### Usage Example

```go
// Store a value in context
c.Set("user_id", 123)

// Later, retrieve it with type safety
userID := c.GetInt("user_id") // Returns 123 as int

// If the wrong type is stored, it returns the zero value
invalidID := c.GetString("user_id") // Returns "" (empty string)

// Example with middleware setting values
func AuthMiddleware(next echoext.HandlerFunc) echoext.HandlerFunc {
    return func(c echoext.Context) error {
        // Set values that can be retrieved with type safety in handlers
        c.Set("user_id", 42)
        c.Set("is_admin", true)
        return next(c)
    }
}

// Handler using the type-safe getters
func MyHandler(c echoext.Context) error {
    userID := c.GetInt("user_id")      // 42
    isAdmin := c.GetBool("is_admin")   // true
    
    // Use the values...
    return c.JSON(http.StatusOK, map[string]interface{}{
        "userId": userID,
        "isAdmin": isAdmin,
    })
}
```


