package routes

import (
	"inquire/now-microservice/handlers"
	"inquire/now-microservice/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes defines all application routes.
func SetupRoutes(router *gin.Engine) {
	// Public routes
	router.GET("/", handlers.HomeHandler)
	router.GET("/submit-form", handlers.SubmitFormHandler)
	router.GET("/geo-check", handlers.GeoCheckHandler)
	router.GET("/close-modal", handlers.CloseModal)
	router.POST("/auth/callback", handlers.AuthCallback)

	// Protected routes group (example)
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// Example protected route
		protected.GET("/profile", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"user":  "John Doe",
				"email": "john.doe@example.com",
			})
		})
	}
}
