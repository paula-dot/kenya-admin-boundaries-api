package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paula-dot/kenya-admin-boundaries-api/internal/service"
)

// SpatialHandler manages HTTP requests for spatial operations
type SpatialHandler struct {
	spatialService *service.SpatialService
}

// NewSpatialHandler initializes the handler
func NewSpatialHandler(s *service.SpatialService) *SpatialHandler {
	return &SpatialHandler{spatialService: s}
}

// IntersectRequest defines the expected JSON payload
type IntersectRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

// HandleIntersect processes the POST request for coordinate intersection
func (h *SpatialHandler) HandleIntersect(c *gin.Context) {
	var req IntersectRequest

	// 1. Input Validation: Bind JSON and enforce required fields
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": "Ensure 'latitude' and 'longitude' are provided as valid numbers.",
		})
		return
	}

	// 2. Execute Business Logic
	resp, err := h.spatialService.GetIntersection(c.Request.Context(), req.Latitude, req.Longitude)
	if err != nil {
		// Catch the specific bounds error from the service layer
		if err.Error() == "invalid coordinates: latitude must be between -90 and 90, and longitude between -180 and 180" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// For all other DB/system errors, return a generic 500 to avoid leaking internals
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process spatial intersection"})
		return
	}

	// 3. Return Successful JSON Response
	c.JSON(http.StatusOK, resp)
}
