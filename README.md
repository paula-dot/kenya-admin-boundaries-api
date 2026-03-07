# Kenya Admin Boundaries API

A production-grade geospatial REST API mapping Kenya's administrative and electoral boundaries down to the constituency and sub-county level.

This project provides high-performance spatial queries, caching, and a clean JSON interface, built with a focus on Clean Architecture and scalable deployment.

## Tech Stack

* **Backend:** Golang 1.25 (Clean Architecture)
* **Web Framework:** Gin
* **Database:** PostgreSQL 15 with PostGIS 3.3
* **Cache:** Redis 7 
* **Database Interactions:** `sqlc` for type-safe SQL, `golang-migrate` for schema definitions
* **Infrastructure:** Docker & Docker Compose

## Repository Structure

The Go REST API is designed using strict separation of concerns:
* `cmd/api/main.go` - The entry point and dependency injection wiring.
* `internal/domain` - Core business models (e.g., County, SubCounty, Constituency).
* `internal/repository/postgres` - `sqlc`-generated database interactions.
* `internal/service` - Business logic and caching interfaces.
* `internal/handler` - HTTP delivery, request binding, and JSON responses.
* `migrations/` - Raw SQL definitions for the PostGIS tables.

## Getting Started

### Prerequisites
* Docker & Docker Compose

### Quick Start

1.  **Start the infrastructure:**
    ```bash
    docker compose up -d
    ```
    This spins up the Go API (Port 18080), PostgreSQL/PostGIS (Port 5432), and Redis (Port 6379), using the overrides specified in `docker-compose.override.yml`.

2.  **Access the API:**
    * Base URL: `http://127.0.0.1:18080/api/v1`
    * Routes: `http://127.0.0.1:18080/__routes`
    * Health check: `http://127.0.0.1:18080/health`


## Core API Endpoints

### 1. Counties
* **GET `/api/v1/counties`**
  Returns a GeoJSON FeatureCollection of all 47 Kenyan counties.
* **GET `/api/v1/counties/:code`**
  Get a specific county boundary as a single GeoJSON Feature. (e.g., `KE047` for Nairobi).

### 2. Constituencies
* **GET `/api/v1/counties/:code/constituencies`**
  Get all constituencies within a specific county, returned as a FeatureCollection.

### 3. Sub-Counties
* **GET `/api/v1/sub-counties`**
  Returns a lightweight JSON array describing all administrative sub-counties.
* **GET `/api/v1/counties/:code/sub-counties`**
  Returns a lightweight JSON array of sub-counties within a specific county.

### 4. Hierarchical Metadata
* **GET `/api/v1/counties/:code/hierarchy`**
  A fast, lightweight endpoint returning the County code/name tightly coupled with an array of its Constituencies, completely omitting the multi-megabyte PostGIS geometries.

### 5. Spatial Intersections
* **POST `/api/v1/spatial/intersect`**
  Submit a Lat/Lng coordinate pair to find exactly which administrative boundaries (County, Constituency, Ward) the point falls inside.
  
  *Example Payload:*
  ```json
  {
      "latitude": -1.286389, 
      "longitude": 36.817223
  }
  ```

## Architecture Notes
* **Clean Architecture:** The backend strictly isolates the domain logic from external frameworks, HTTP delivery, and database implementations.
* **Spatial Indexing:** All boundary geometries use PostGIS `GIST` indexes for highly optimized sub-millisecond spatial ST_Intersects operations.
* **Caching Strategy:** Expensive `ST_AsGeoJSON` database queries are cached in Redis to guarantee sub-millisecond response times for frequently requested boundaries.
* **Service Interfaces:** Handlers depend on strictly defined Service interfaces, making unit-testing and mock generation simple.
