package routes

import (
	"go-backend/internal/handlers"
	"go-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the routes for the given gin.Engine
// It registers the /health, /login, /users, and /logout endpoints.
func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")

	api.GET("/health", handlers.HealthCheck)
	api.POST("/signup", handlers.CreateUser)
	api.POST("/login", handlers.Login)
	api.POST("/refresh", handlers.RefreshToken)

	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	// protected.GET("/user", handlers.GetUser)
	protected.POST("/logout", handlers.Logout)
	protected.POST("/logout-all", handlers.LogoutAll)
}