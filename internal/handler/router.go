package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/paula-dot/kenya-admin-boundaries-api/internal/domain"
	"github.com/paula-dot/kenya-admin-boundaries-api/internal/service"
)

// SetupRouter wires service into HTTP routes and returns a *gin.Engine.
// Routes implemented (from README):
// GET /api/v1/counties - list all counties (FeatureCollection)
// GET /api/v1/counties/:slug - single county (FeatureCollection)
// GET /api/v1/counties/:slug/constituencies - constituencies in a county (FeatureCollection)
// GET /api/v1/constituencies/:slug/wards - wards in a constituency (FeatureCollection)
// POST /api/v1/spatial/intersect - submit { "lat": <float>, "lng": <float> } and return matching ward/constituency/county (FeatureCollection)
func SetupRouter(svc *service.CountyService) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		// List all counties using the service method that returns a FeatureCollection
		v1.GET("/counties", func(c *gin.Context) {
			ctx := c.Request.Context()
			fc, err := svc.ListCountiesAsFeatureCollection(ctx)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			out, _ := json.Marshal(fc)
			c.Data(http.StatusOK, "application/geo+json", out)
		})

		// The County-by-slug, constituencies, wards and spatial intersect routes depend
		// on additional service methods that may not be implemented on the core
		// CountyService yet. We perform a runtime type-assertion and return 501
		// Not Implemented when those helpers are not available.

		// GET /api/v1/counties/:slug
		v1.GET("/counties/:slug", func(c *gin.Context) {
			slug := c.Param("slug")
			ctx := c.Request.Context()

			// optional interface for slug lookup
			type slugSvc interface {
				GetCountyBySlug(ctx context.Context, slug string) (*domain.County, error)
			}

			if s, ok := interface{}(svc).(slugSvc); ok {
				cnt, err := s.GetCountyBySlug(ctx, slug)
				if err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
					return
				}
				fc := buildFeatureCollectionFromCounties([]*domain.County{cnt})
				out, _ := json.Marshal(fc)
				c.Data(http.StatusOK, "application/geo+json", out)
				return
			}

			c.JSON(http.StatusNotImplemented, gin.H{"error": "GetCountyBySlug not implemented in service"})
		})

		// GET /api/v1/counties/:slug/constituencies
		v1.GET("/counties/:slug/constituencies", func(c *gin.Context) {
			slug := c.Param("slug")
			ctx := c.Request.Context()

			type constSvc interface {
				ListConstituenciesByCountySlug(ctx context.Context, slug string) ([]struct {
					ID       int32
					Slug     string
					Name     string
					Geometry []byte
				}, error)
			}

			if s, ok := interface{}(svc).(constSvc); ok {
				list, err := s.ListConstituenciesByCountySlug(ctx, slug)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				fc := buildFeatureCollectionFromConstituencies(list)
				out, _ := json.Marshal(fc)
				c.Data(http.StatusOK, "application/geo+json", out)
				return
			}

			c.JSON(http.StatusNotImplemented, gin.H{"error": "ListConstituenciesByCountySlug not implemented in service"})
		})

		// GET /api/v1/constituencies/:slug/wards
		v1.GET("/constituencies/:slug/wards", func(c *gin.Context) {
			slug := c.Param("slug")
			ctx := c.Request.Context()

			type wardSvc interface {
				ListWardsByConstituencySlug(ctx context.Context, slug string) ([]struct {
					ID       int32
					Slug     string
					Name     string
					Geometry []byte
				}, error)
			}

			if s, ok := interface{}(svc).(wardSvc); ok {
				list, err := s.ListWardsByConstituencySlug(ctx, slug)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				fc := buildFeatureCollectionFromWards(list)
				out, _ := json.Marshal(fc)
				c.Data(http.StatusOK, "application/geo+json", out)
				return
			}

			c.JSON(http.StatusNotImplemented, gin.H{"error": "ListWardsByConstituencySlug not implemented in service"})
		})

		// POST /api/v1/spatial/intersect
		v1.POST("/spatial/intersect", func(c *gin.Context) {
			var payload struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			}
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
				return
			}
			ctx := c.Request.Context()

			type spatialSvc interface {
				SpatialIntersect(ctx context.Context, lat, lng float64) (struct {
					Ward         *domain.Ward
					Constituency *domain.Constituency
					County       *domain.County
				}, error)
			}

			if s, ok := interface{}(svc).(spatialSvc); ok {
				res, err := s.SpatialIntersect(ctx, payload.Lat, payload.Lng)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				// Build a FeatureCollection containing up to Ward, Constituency, County in that order
				var features []interface{}
				if res.Ward != nil {
					// domain.Ward exposes GeoJSON as a string field; convert to []byte for the helper
					features = append(features, buildFeature([]byte(res.Ward.GeoJSON), map[string]interface{}{
						"type": "ward",
						"id":   res.Ward.ID,
						"name": res.Ward.Name,
					}))
				}
				if res.Constituency != nil {
					features = append(features, buildFeature([]byte(res.Constituency.GeoJSON), map[string]interface{}{
						"type": "constituency",
						"id":   res.Constituency.ID,
						"name": res.Constituency.Name,
					}))
				}
				if res.County != nil {
					features = append(features, buildFeature(res.County.Geometry, map[string]interface{}{
						"type": "county",
						"id":   res.County.ID,
						"name": res.County.Name,
					}))
				}

				fc := map[string]interface{}{
					"type":     "FeatureCollection",
					"features": features,
				}
				out, _ := json.Marshal(fc)
				c.Data(http.StatusOK, "application/geo+json", out)
				return
			}

			c.JSON(http.StatusNotImplemented, gin.H{"error": "SpatialIntersect not implemented in service"})
		})
	}

	return r
}

// helpers to construct GeoJSON features/collections
func buildFeatureCollectionFromCounties(counties []*domain.County) map[string]interface{} {
	features := make([]interface{}, 0, len(counties))
	for _, c := range counties {
		props := map[string]interface{}{
			"id":   c.ID,
			"code": c.Code,
			"name": c.Name,
		}
		features = append(features, buildFeature(c.Geometry, props))
	}
	return map[string]interface{}{
		"type":     "FeatureCollection",
		"features": features,
	}
}

func buildFeatureCollectionFromConstituencies(list []struct {
	ID       int32
	Slug     string
	Name     string
	Geometry []byte
}) map[string]interface{} {
	features := make([]interface{}, 0, len(list))
	for _, it := range list {
		props := map[string]interface{}{
			"id":   it.ID,
			"slug": it.Slug,
			"name": it.Name,
		}
		features = append(features, buildFeature(it.Geometry, props))
	}
	return map[string]interface{}{
		"type":     "FeatureCollection",
		"features": features,
	}
}

func buildFeatureCollectionFromWards(list []struct {
	ID       int32
	Slug     string
	Name     string
	Geometry []byte
}) map[string]interface{} {
	features := make([]interface{}, 0, len(list))
	for _, it := range list {
		props := map[string]interface{}{
			"id":   it.ID,
			"slug": it.Slug,
			"name": it.Name,
		}
		features = append(features, buildFeature(it.Geometry, props))
	}
	return map[string]interface{}{
		"type":     "FeatureCollection",
		"features": features,
	}
}

func buildFeature(geometry []byte, properties map[string]interface{}) map[string]interface{} {
	// geometry is expected to be raw JSON bytes representing a geometry object
	var geom json.RawMessage
	if len(geometry) > 0 {
		geom = json.RawMessage(geometry)
	} else {
		geom = json.RawMessage([]byte("null"))
	}
	return map[string]interface{}{
		"type":       "Feature",
		"geometry":   geom,
		"properties": properties,
	}
}

// Note: we add a small helper on domain types usage where needed. If domain.Ward
// does not expose a GeoJSON byte getter, we attempt to use GeoJSON string fields.
// This file intentionally avoids making the handler depend on many service
// implementation details; runtime assertions are used to keep compilation safe.
