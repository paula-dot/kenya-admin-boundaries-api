-- name: CreateCounty :one
-- Inserts a new county and converts the incoming GeoJSON payload into a PostGIS geometry.
INSERT INTO counties (code, name, geom)
VALUES (
        $1,
        $2,
        ST_Multi(ST_GeomFromGeoJSON($3::text))
       )
RETURNING id, code, name, created_at

-- name: GetCountyByID :one
-- Fetches a specific county and automatically formats the geometry as valid GeoJSON.
SELECT
    id,
    code,
    name,
    ST_AsGeoJSON(geom)::jsonb AS geojson,
    created_at,
    updated_at
FROM counties
WHERE id = $1 LIMIT 1;

-- name: ListCounties :many
-- Retrieves a list of all counties.
SELECT
    id,
    code,
    name,
    ST_AsGeoJSON(geom)::jsonb AS geojson
FROM counties;
ORDER BY name;