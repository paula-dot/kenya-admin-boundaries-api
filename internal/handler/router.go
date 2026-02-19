package handler

import (
	"github/gin-gonic/gin"
)

// SetupRouter configures the Gin engine and registers all API routes.
func SetupRouter(CountyHandler *CountyHandler) *gin.Engine {
	// Use gin.New() instead of Default() if you want to explicitly define your middlewares later
	router := gin.New()

	// Global Middlewares can be added here (e.g., CORS, Rate Limiting)
	// router.Use(CORSMiddleware())

	// Health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "API is operational"})
	})

	// API Versioning Group
	v1 := router.Group("/api/v1")
	{
		// County routes
		counties := v1.Group("/counties")
		{
			counties.GET("", countyHandler, ListCounties)
			counties.GET("/:id", countyHandler.GetCounty)
		}
	}

	return router
}
