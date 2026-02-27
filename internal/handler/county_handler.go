package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

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
	code := c.Param("slug")
	feature, err := h.countyService.GetCountyAsFeature(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve county"})
		return
	}

	c.JSON(http.StatusOK, feature)
}

// ListCounties retrieves all counties as a GeoJSON FeatureCollection.
func (h *CountyHandler) ListCounties(c *gin.Context) {
	collection, err := h.countyService.ListCountiesAsFeatureCollection(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve counties"})
		return
	}

	c.JSON(http.StatusOK, collection)
}
