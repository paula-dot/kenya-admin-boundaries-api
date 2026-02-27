package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/paula-dot/kenya-admin-boundaries-api/internal/dto"
	"github.com/paula-dot/kenya-admin-boundaries-api/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type APIHandler struct {
	Repo *repository.SpatialRepository
}

// LocationResponse defines your JSON output
type LocationResponse struct {
	CountyCode       string `json:"county_code"`
	CountyName       string `json:"county_name"`
	ConstituencyCode string `json:"constituency_code,omitempty"`
	ConstituencyName string `json:"constituency_name,omitempty"`
}

// GetLocationByPoint handles the reverse-geocoding API request
func GetLocationByPoint(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Extract parameters (?lat=-0.514&lng=35.269)
		latStr := c.Query("lat")
		lngStr := c.Query("lng")

		lat, err := strconv.ParseFloat(latStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Valid latitude is required"})
			return
		}

		lng, err := strconv.ParseFloat(lngStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Valid longitude is required"})
			return
		}

		query := `
         	SELECT c.county_code, c.county_name, con.constituency_code, con.constituency_name
			FROM constituencies con
         	JOIN counties c ON con.county_code = c.county_code
         	WHERE ST_Intersects(con.geom, ST_SetSRID(ST_MakePoint($1, $2), 4326));
        `

		var res LocationResponse
		err = db.QueryRow(context.Background(), query, lng, lat).Scan(
			&res.CountyCode, &res.CountyName, &res.ConstituencyCode, &res.ConstituencyName,
		)

		// 3. Handle Errors
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Coordinates fall outside defined administrative boundaries"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal database error"})
			return
		}

		// 4. Success Response
		c.JSON(http.StatusOK, res)
	}
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
	var req dto.PointRequest

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
