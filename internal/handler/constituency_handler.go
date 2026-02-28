package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/paula-dot/kenya-admin-boundaries-api/internal/service"
)

// ConstituencyHandler handles HTTP requests for constituency data.
type ConstituencyHandler struct {
	service *service.ConstituencyService
}

// NewConstituencyHandler creates a new handler.
func NewConstituencyHandler(svc *service.ConstituencyService) *ConstituencyHandler {
	return &ConstituencyHandler{service: svc}
}

// ListByCounty responds with a FeatureCollection of constituencies for a county slug.
func (h *ConstituencyHandler) ListByCounty(c *gin.Context) {
	slug := c.Param("slug")
	ctx := c.Request.Context()
	list, err := h.service.ListConstituenciesByCountySlug(ctx, slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fc := buildFeatureCollectionFromConstituencies(list)
	out, _ := json.Marshal(fc)
	c.Data(http.StatusOK, "application/geo+json", out)
}
