package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/paula-dot/kenya-admin-boundaries-api/internal/dto/request"
	"github.com/paula-dot/kenya-admin-boundaries-api/internal/repository"
)

type APIHandler struct {
	Repo *repository.SpatialRepository
}

// --- Handlers ---
func listCounties(c *gin.Context) {
	// Future: Fetch from DB and return as GeoJSON
	c.JSON(http.StatusOK, gin.H{"message": "List all counties as GeoJSON FeatureCollection"})
}

func (h *APIHandler) GetCountyBySlug(c *gin.Context) {
	slug := c.Param("slug")

	// Pass the request context down to the database for graceful timeouts
	feature, err := h.Repo.GetCountyBySlug(c.Request.Context(), slug)
	if err != nil {
		// Log the actual error internally, return a generic 404 to the user
		c.JSON(http.StatusNotFound, gin.H{"error": "County not found"})
		return
	}

	// Return the perfectly formatted GeoJSON
	c.JSON(http.StatusOK, feature)
}

func getConstituenciesByCounty(c *gin.Context) {
	slug := c.Param("slug")
	c.JSON(http.StatusOK, gin.H{"message": "Get constituencies for county: " + slug})
}

func getWardsByConstituency(c *gin.Context) {
	slug := c.Param("slug") // e.g., "kapseret"
	c.JSON(http.StatusOK, gin.H{"message": "Get wards for constituency: " + slug})
}

func (h *APIHandler) CheckIntersection(c *gin.Context) {
	var req request.PointRequest

	// BindJSON validates the incoming payload against our struct tags (min/max limits)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid coordinates provided", "details": err.Error()})
		return
	}

	result, err := h.Repo.GetLocationByPoint(c.Request.Context(), req.Longitude, req.Latitude)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coordinates do not fall within known boundaries"})
		return
	}

	c.JSON(http.StatusOK, result)
}
