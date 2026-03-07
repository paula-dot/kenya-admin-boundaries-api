package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/paula-dot/kenya-admin-boundaries-api/internal/domain"
	"github.com/paula-dot/kenya-admin-boundaries-api/internal/service"
	"github.com/paula-dot/kenya-admin-boundaries-api/pkg/geojson"
)

// AppServices is a small container for wiring explicit handlers in SetupRouter.
// When passed into SetupRouter callers can avoid runtime assertions and the
// router will register concrete handlers using these service instances.
type AppServices struct {
	County       *service.CountyService
	Constituency *service.ConstituencyService
	Spatial      *service.SpatialService
}

// SetupRouter wires service into HTTP routes and returns a *gin.Engine.
// Accepts optional middleware functions that will be applied to the /api/v1
// route group. This allows callers to attach rate-limiters or auth middleware
// without modifying the internal router implementation.
func SetupRouter(svc interface{}, v1Middleware ...gin.HandlerFunc) *gin.Engine {
	r := gin.Default()

	// Debug middleware: log raw request info to help diagnose unexpected path rewrites
	// (kept minimal and only enabled in development). It logs RequestURI, URL.Path
	// and a few proxy headers which often cause path-prefixing issues.
	r.Use(func(c *gin.Context) {
		log.Printf("[REQ] Method=%s RequestURI=%s URL.Path=%s Host=%s Referer=%s X-Forwarded-Host=%s X-Forwarded-Proto=%s\n",
			c.Request.Method,
			c.Request.RequestURI,
			c.Request.URL.Path,
			c.Request.Host,
			c.GetHeader("Referer"),
			c.GetHeader("X-Forwarded-Host"),
			c.GetHeader("X-Forwarded-Proto"),
		)
		c.Next()
	})

	// create the /api/v1 group and apply any provided middleware
	v1 := r.Group("/api/v1", v1Middleware...)
	{
		// List all counties using the service method that returns a FeatureCollection
		v1.GET("/counties", func(c *gin.Context) {
			ctx := c.Request.Context()

			// runtime assertion for ListCounties
			type listSvc interface {
				ListCountiesAsFeatureCollection(ctx context.Context) (*geojson.FeatureCollection, error)
			}

			// Fallback: if svc implements the concrete CountyService, call its method directly
			if s, ok := svc.(*service.CountyService); ok {
				fc, err := s.ListCountiesAsFeatureCollection(ctx)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				out, _ := json.Marshal(fc)
				c.Data(http.StatusOK, "application/geo+json", out)
				return
			}

			// If svc implements the interface, use it (keeps tests flexible)
			if s, ok := svc.(listSvc); ok {
				fc, err := s.ListCountiesAsFeatureCollection(context.Background())
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				out, _ := json.Marshal(fc)
				c.Data(http.StatusOK, "application/geo+json", out)
				return
			}

			c.JSON(http.StatusNotImplemented, gin.H{"error": "ListCountiesAsFeatureCollection not implemented in service"})
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

			if s, ok := svc.(slugSvc); ok {
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

		// Register constituencies by county route. Prefer explicit handler wiring
		// when caller passed an *AppServices. Fall back to the existing runtime
		// assertion behaviour for test flexibility.
		if svcApp, ok := svc.(*AppServices); ok && svcApp.Constituency != nil {
			ch := NewConstituencyHandler(svcApp.Constituency)
			v1.GET("/counties/:slug/constituencies", ch.ListByCounty)
		} else {
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

				if s, ok := svc.(constSvc); ok {
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
		}

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

			if s, ok := svc.(wardSvc); ok {
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
		if svcApp, ok := svc.(*AppServices); ok && svcApp.Spatial != nil {
			spatialHdlr := NewSpatialHandler(svcApp.Spatial)
			v1.POST("/spatial/intersect", spatialHdlr.HandleIntersect)
		} else {
			// Fallback: runtime assertion for test flexibility
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
					SpatialIntersect(ctx context.Context, lat, lng float64) (service.SpatialResult, error)
				}

				if s, ok := svc.(spatialSvc); ok {
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

		// GET /api/v1/location (query params lat,lng) - runtime assertion against spatialSvc
		v1.GET("/location", func(c *gin.Context) {
			latStr := c.Query("lat")
			lngStr := c.Query("lng")
			if latStr == "" || lngStr == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "lat and lng query parameters are required"})
				return
			}
			lat, err1 := strconv.ParseFloat(latStr, 64)
			lng, err2 := strconv.ParseFloat(lngStr, 64)
			if err1 != nil || err2 != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "lat and lng must be valid floats"})
				return
			}
			ctx := c.Request.Context()

			type spatialSvc interface {
				SpatialIntersect(ctx context.Context, lat, lng float64) (service.SpatialResult, error)
			}

			if s, ok := svc.(spatialSvc); ok {
				res, err := s.SpatialIntersect(ctx, lat, lng)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				var features []interface{}
				if res.Ward != nil {
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
