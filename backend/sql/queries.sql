-- name: CreateCounty :one
-- Inserts a new county and converts the incoming GeoJSON payload into a PostGIS geometry.
INSERT INTO counties (county_code, county_name, geom)
VALUES (
        $1,
        $2,
        ST_Multi(ST_GeomFromGeoJSON($3::text))
       )
RETURNING county_code AS code, county_name AS name;

-- name: GetCountyByCode :one
-- Fetches a specific county by its official code and automatically formats the geometry as valid GeoJSON.
SELECT
    county_code AS code,
    county_name AS name,
    ST_AsGeoJSON(geom)::jsonb AS geojson
FROM counties
WHERE county_code = $1
LIMIT 1;

-- name: ListCounties :many
-- Retrieves a list of all counties. Uses code/name and numeric ordering.
SELECT
    county_code AS id,
    county_name AS name,
    ST_AsGeoJSON(geom)::jsonb AS geojson
FROM counties
ORDER BY county_code ASC;

-- name: CreateConstituency :one
-- Inserts a new constituency and converts the incoming GeoJSON payload into a PostGIS geometry.
INSERT INTO constituencies (constituency_code, constituency_name, county_code, geom)
VALUES (
        $1,
        $2,
        $3,
        ST_Multi(ST_GeomFromGeoJSON($4::text))
       )
RETURNING constituency_code AS code, constituency_name AS name;

-- name: GetConstituencyByCode :one
-- Fetches a specific constituency by its official code.
SELECT
    constituency_code AS id,
    constituency_name AS name,
    county_code,
    ST_AsGeoJSON(geom)::jsonb AS geojson
FROM constituencies
WHERE constituency_code = $1
LIMIT 1;

-- name: ListConstituencies :many
-- Retrieves a list of all constituencies.
SELECT
    constituency_code AS id,
    constituency_name AS name,
    county_code,
    ST_AsGeoJSON(geom)::jsonb AS geojson
FROM constituencies
ORDER BY constituency_code ASC;

-- name: ListConstituenciesByCounty :many
-- Retrieves all constituencies belonging to a specific county (Perfect for your nested route!)
SELECT
    constituency_code AS id,
    constituency_name AS name,
    county_code,
    ST_AsGeoJSON(geom)::jsonb AS geojson
FROM constituencies
WHERE county_code = $1
ORDER BY constituency_code ASC;

-- name: GetIntersectingBoundary :one
SELECT
    c.county_code,
    c.county_name,
    co.constituency_code,
    co.constituency_name
FROM constituencies co
JOIN counties c ON co.county_code = c.county_code
WHERE ST_Intersects(
      co.geom,
      ST_SetSRID(ST_MakePoint(@longitude::float, @latitude::float), 4326)
      )
LIMIT 1;
