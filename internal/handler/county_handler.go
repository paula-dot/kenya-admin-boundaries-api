package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	
	"github.com/paula-dot/kenya-admin-boundaries-api/internal/service"
)

// CountyHandler handles HTTP requests for county data.
type CountyHandler struct {
	countyService *service.CountyService
}

// NewCountyHandler creates a new instance of CountyHandler.
func NewCountyHandler(countyService *service.CountyService) *CountyHandler {
	return &CountyHandler{
		countyService: countyService,
	}
}

// GetCounty retrieves a single county as a GeoJSON Feature.
func (h *CountyHandler) GetCounty(c *gin.Context) {
	// 1. Input Validation: Ensure the ID is a valid UUID to prevent SQL injection or bad queries
	idParam := c.Param("id")
	countyID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid county ID format. Must be a UUID."})
		return
	}

	// 2. Call the Service Layer
	feature, err := h.countyService.GetCountyAsFeature(c.Request.Context(), countyID)
	if err != nil {
		// In a production app, you'd check if the error is a "not found" vs "internal server error"
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve county"})
		return
	}

	// 3. Return the payload. Gin will perfectly serialize the json.RawMessage inside the Feature.
	c.JSON(http.StatusOK, feature)
}

// ListCounties retrieves all counties as a GeoJSON FeatureCollection.
func (h *CountyHandler) ListCounties(c *gin.Context) {
	collection, err := h.countyService.ListCountiesAsFeatureCollection(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginH{"error": "Failed to retrieve counties"})
		return
	}

	c.JSON(http.StatusOK, collection)
}
