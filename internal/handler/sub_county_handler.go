package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paula-dot/kenya-admin-boundaries-api/internal/service"
)

// SubCountyHandler provides HTTP handlers for sub-county endpoints.
type SubCountyHandler struct {
	svc *service.SubCountyService
}

// NewSubCountyHandler creates a new handler.
func NewSubCountyHandler(svc *service.SubCountyService) *SubCountyHandler {
	return &SubCountyHandler{svc: svc}
}

// ListAll returns the list of all sub-counties.
func (h *SubCountyHandler) ListAll(c *gin.Context) {
	results, err := h.svc.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sub-counties"})
		return
	}

	// Make sure an empty list becomes [] rather than null in JSON
	if len(results) == 0 {
		c.JSON(http.StatusOK, gin.H{"sub_counties": []interface{}{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sub_counties": results})
}

// ListByCounty returns sub-counties nested inside a particular county code.
func (h *SubCountyHandler) ListByCounty(c *gin.Context) {
	code := c.Param("slug") // we reuse the param "slug" to mean code in these dynamic routes
	results, err := h.svc.GetByCountyCode(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sub-counties for county"})
		return
	}

	if len(results) == 0 {
		c.JSON(http.StatusOK, gin.H{"sub_counties": []interface{}{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sub_counties": results})
}
