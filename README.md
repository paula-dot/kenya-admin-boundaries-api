#  Kenya County, Constituency & Ward API

A production-grade geospatial API and visualization platform specifically mapping Kenya's administrative and electoral boundaries down to the ward level.

This project provides high-performance spatial queries, caching, and a responsive web interface, built with a focus on Clean Architecture and scalable deployment.

##  Tech Stack

* **Backend:** Golang (Clean Architecture)
* **Database:** PostgreSQL with PostGIS extension
* **Cache:** Redis
* **Frontend:** React (TypeScript) with Leaflet.js
* **ETL Pipeline:** Python (GeoPandas, SQLAlchemy)
* **Infrastructure:** Docker & Docker Compose

##  Project Structure

This repository is organized as a monorepo containing three main services:

* `/backend` - The Go REST API, implementing strict separation of concerns (Domain, Service, Repository, Handler).
* `/frontend` - The React application for map visualization and boundary searching.
* `/data-pipeline` - Python scripts for extracting `.gpkg` or `.shp` boundary data and loading it into PostGIS.

##  Getting Started

###  Prerequisites

Ensure you have the following installed on your machine:
* [Docker](https://docs.docker.com/get-docker/)
* [Docker Compose](https://docs.docker.com/compose/install/)
* Git

###  Quick Start

1.  **Clone the repository:**
    ```bash
    git clone [https://github.com/your-org/kenya-county-constituency-ward-api.git](https://github.com/your-org/kenya-county-constituency-ward-api.git)
    cd kenya-county-constituency-ward-api
    ```

2.  **Start the infrastructure (Postgres/PostGIS, Redis, API, Frontend):**
    ```bash
    docker-compose up -d --build
    ```

3.  **Run Database Migrations:**
    ```bash
    docker-compose exec backend migrate -path ./migrations -database "$DATABASE_URL" up
    ```

4.  **Load Initial Boundary Data:**
    Place your raw GeoPackage files in `data-pipeline/raw_data/` and run the ETL script:
    ```bash
    cd data-pipeline
    python -m venv venv
    source venv/bin/activate
    pip install -r requirements.txt
    python scripts/load_data.py
    ```

5.  **Access the Application:**
    * **Frontend UI:** `http://localhost:3000`
    * **API Base URL:** `http://localhost:8080/api/v1`
    * **API Documentation (Swagger):** `http://localhost:8080/swagger/index.html`

##  Core API Endpoints

All spatial responses are returned as standard `GeoJSON FeatureCollection` objects.

* `GET /api/v1/counties` - List all counties.
* `GET /api/v1/counties/:slug` - Get specific county boundaries.
* `GET /api/v1/counties/:slug/constituencies` - Get constituencies within a specific county.
* `GET /api/v1/constituencies/:slug/wards` - Get wards within a specific constituency.
* `POST /api/v1/spatial/intersect` - Submit a Point (Lat/Lng) to find exactly which Ward, Constituency, and County it falls inside.

##  Architecture & Best Practices

* **Clean Architecture:** The backend strictly isolates the domain logic from external frameworks, HTTP delivery, and database implementations.
* **Spatial Indexing:** All boundary geometries utilize PostGIS `GIST` indexes for highly optimized spatial intersections.
* **Caching Strategy:** Expensive `ST_AsGeoJSON` database queries are cached in Redis to guarantee sub-millisecond response times for frequently requested boundaries.