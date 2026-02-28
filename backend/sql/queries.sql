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
