package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/BacoFoods/echoext"
)

// User represents a simple user model
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// UserService is a simple mock service for user operations
type UserService struct {
	users map[int]User
}

// NewUserService creates a new instance of UserService
func NewUserService() *UserService {
	return &UserService{
		users: map[int]User{
			1: {ID: 1, Name: "John Doe", Email: "john@example.com", CreatedAt: time.Now()},
			2: {ID: 2, Name: "Jane Smith", Email: "jane@example.com", CreatedAt: time.Now()},
		},
	}
}

// Handler is our application handler containing all dependencies
type Handler struct {
	userService *UserService
}

// NewHandler creates a new instance of Handler
func NewHandler() *Handler {
	return &Handler{
		userService: NewUserService(),
	}
}

// GetUsers returns all users
func (h *Handler) GetUsers(c echoext.Context) error {
	users := make([]User, 0, len(h.userService.users))
	for _, user := range h.userService.users {
		users = append(users, user)
	}
	return c.JSON(http.StatusOK, users)
}

// GetUser returns a specific user by ID
func (h *Handler) GetUser(c echoext.Context) error {
	id := c.GetInt("id")
	if id == 0 {
		// Try to parse from param
		paramID := c.Param("id")
		// In a real app, you should handle the error properly
		id, _ = parseInt(paramID)
	}

	user, exists := h.userService.users[id]
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, user)
}

// CreateUser creates a new user
func (h *Handler) CreateUser(c echoext.Context) error {
	var user User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user data"})
	}

	// Simple validation
	if user.Name == "" || user.Email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Name and email are required"})
	}

	// In a real app, you'd generate a unique ID
	user.ID = len(h.userService.users) + 1
	user.CreatedAt = time.Now()

	h.userService.users[user.ID] = user

	return c.JSON(http.StatusCreated, user)
}

// UpdateUser updates an existing user
func (h *Handler) UpdateUser(c echoext.Context) error {
	id := c.GetInt("id")
	if id == 0 {
		// Try to parse from param
		paramID := c.Param("id")
		// In a real app, you should handle the error properly
		id, _ = parseInt(paramID)
	}

	user, exists := h.userService.users[id]
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	var updatedUser User
	if err := c.Bind(&updatedUser); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user data"})
	}

	// Only update non-empty fields
	if updatedUser.Name != "" {
		user.Name = updatedUser.Name
	}
	if updatedUser.Email != "" {
		user.Email = updatedUser.Email
	}

	h.userService.users[id] = user

	return c.JSON(http.StatusOK, user)
}

// DeleteUser deletes a user
func (h *Handler) DeleteUser(c echoext.Context) error {
	id := c.GetInt("id")
	if id == 0 {
		// Try to parse from param
		paramID := c.Param("id")
		// In a real app, you should handle the error properly
		id, _ = parseInt(paramID)
	}

	_, exists := h.userService.users[id]
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	delete(h.userService.users, id)

	return c.NoContent(http.StatusNoContent)
}

// LoggingMiddleware is a simple middleware to log requests
func LoggingMiddleware(next echoext.HandlerFunc) echoext.HandlerFunc {
	return func(c echoext.Context) error {
		start := time.Now()

		// Call the next handler
		err := next(c)

		// Log the request
		method := c.Request().Method
		path := c.Request().URL.Path
		duration := time.Since(start)

		c.Logger().Infof("[%s] %s - %v", method, path, duration)

		return err
	}
}

// AuthMiddleware is a simple middleware to check for authentication
func AuthMiddleware(next echoext.HandlerFunc) echoext.HandlerFunc {
	return func(c echoext.Context) error {
		// In a real app, you would check for a valid token
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authorization required"})
		}

		// For demo purposes, let's set a user ID
		c.Set("id", 1)

		return next(c)
	}
}

// parseInt is a helper function to convert string to int
func parseInt(s string) (int, error) {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}

// StartServer starts the example server
func StartServer() {
	// Create a server with custom configuration
	config := echoext.ServerConfig{
		PathPrefix:       "/api/v1",
		Host:             "localhost",
		Port:             3000,
		HealthcheckPath:  "/health",
		ExtraCORSHeaders: []string{"X-Auth-Token", "X-Custom-Header"},
	}

	server := echoext.New(config)

	handler := NewHandler()

	// Create a users group
	server.Group("/users", func(g *echoext.Group) {
		// Register routes with our custom handlers and middleware
		g.GET("", handler.GetUsers, LoggingMiddleware)
		g.GET("/:id", handler.GetUser, LoggingMiddleware)
		g.POST("", handler.CreateUser, LoggingMiddleware, AuthMiddleware)
		g.PUT("/:id", handler.UpdateUser, LoggingMiddleware, AuthMiddleware)
		g.DELETE("/:id", handler.DeleteUser, LoggingMiddleware, AuthMiddleware)
	})

	// Start the server
	server.Start()
}
