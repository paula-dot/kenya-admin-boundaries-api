# Kenya Admin Boundaries API

A production-grade geospatial REST API mapping Kenya's administrative and electoral boundaries down to the constituency level.

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


## Comprehensive API Documentation

### 1. Counties

**`GET /api/v1/counties`**
- **Description:** Returns a GeoJSON FeatureCollection of all 47 Kenyan counties with their polygons and metadata.

**`GET /api/v1/counties/:slug`**
- **Description:** Get a specific county boundary as a single GeoJSON Feature.
- **Parameters:**
  - `slug` (string): The county code (e.g., `KE047` for Nairobi).

**`GET /api/v1/counties/:slug/hierarchy`**
- **Description:** A fast, lightweight endpoint returning the County code and name tightly coupled with an array of its Constituencies, completely omitting the heavy PostGIS geometries.
- **Parameters:**
  - `slug` (string): The county code (e.g., `KE001`).

---

### 2. Constituencies

**`GET /api/v1/constituencies`**
- **Description:** Returns a GeoJSON FeatureCollection of all electoral constituencies in Kenya.

**`GET /api/v1/counties/:slug/constituencies`**
- **Description:** Get all constituencies within a specific county, returned as a FeatureCollection.
- **Parameters:**
  - `slug` (string): The county code (e.g., `KE047` for Nairobi).

---

### 3. Sub-Counties

**`GET /api/v1/sub-counties`**
- **Description:** Returns a lightweight JSON array describing all administrative sub-counties.

**`GET /api/v1/counties/:slug/sub-counties`**
- **Description:** Returns a lightweight JSON array of sub-counties within a specific county.
- **Parameters:**
  - `slug` (string): The county code (e.g., `KE047` for Nairobi).

---

### 4. Wards

**`GET /api/v1/wards`**
- **Description:** Returns a paginated JSON list of all political wards across Kenya, including their parent constituency and county metadata.
- **Parameters:**
  - `page` (number): The page number (default 1).
  - `limit` (number): The number of results per page (default 50).

---

## Interactive Map Playground

While a custom frontend application isn't necessary for an API (and industry standard heavily favors generating docs directly or putting them in the `README`), this project includes a helpful interactive geospatial playground to instantly visualize the PostGIS spatial data on an actual map.

To run the playground:
```bash
cd kenya-admin-ui
npm install
npm run dev
```
Navigate to `http://localhost:5173/map` to browse and interact with the polygon bounds.

## Architecture Notes
* **Clean Architecture:** The backend strictly isolates the domain logic from external frameworks, HTTP delivery, and database implementations.
* **Spatial Indexing:** All boundary geometries use PostGIS `GIST` indexes for highly optimized sub-millisecond spatial ST_Intersects operations.
* **Caching Strategy:** Expensive `ST_AsGeoJSON` database queries are cached in Redis to guarantee sub-millisecond response times for frequently requested boundaries.
* **Service Interfaces:** Handlers depend on strictly defined Service interfaces, making unit-testing and mock generation simple.
