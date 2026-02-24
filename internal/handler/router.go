package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/paula-dot/kenya-admin-boundaries-api/internal/service"
)

// SetupRouter configures the Gin engine and registers all API routes.
func SetupRouter(svc *service.CountyService) *gin.Engine {
	// Create handlers with service dependencies
	countyHandler := NewCountyHandler(svc)

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
			counties.GET("", countyHandler.ListCounties)
			counties.GET(":id", countyHandler.GetCounty)
		}
	}

	return router
}
