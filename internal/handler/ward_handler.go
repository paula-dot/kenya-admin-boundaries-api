package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/paula-dot/kenya-admin-boundaries-api/internal/repository/postgres"
	"github.com/paula-dot/kenya-admin-boundaries-api/internal/service"
)

// WardHandler serves ward endpoints
type WardHandler struct {
	service *service.WardService
}

// NewWardHandler creates a new handler
func NewWardHandler(svc *service.WardService) *WardHandler {
	return &WardHandler{service: svc}
}

// ListAll responds with a paginated list of wards
func (h *WardHandler) ListAll(c *gin.Context) {
	ctx := c.Request.Context()

	// Default pagination values
	var limit int32 = 50
	var page int32 = 1

	if limitStr := c.Query("limit"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 && parsed <= 500 {
			limit = int32(parsed)
		}
	}

	if pageStr := c.Query("page"); pageStr != "" {
		if parsed, err := strconv.Atoi(pageStr); err == nil && parsed > 0 {
			page = int32(parsed)
		}
	}

	offset := (page - 1) * limit

	wards, err := h.service.ListWards(ctx, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if wards == nil {
		// return empty array instead of null
		wards = make([]postgres.ListWardsRow, 0) // fallback to avoid nil slice in JSON output
	}

	c.JSON(http.StatusOK, gin.H{
		"page":  page,
		"limit": limit,
		"data":  wards,
	})
}
