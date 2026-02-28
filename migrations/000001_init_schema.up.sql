-- Enable PostGIS extension
CREATE EXTENSION IF NOT EXISTS postgis;

-- 1. Counties (Geometry included)
CREATE TABLE counties (
    county_code VARCHAR(50) PRIMARY KEY, -- Official government code (new column used by sqlc queries)
    county_name VARCHAR(100) NOT NULL,   -- New column used by sqlc queries
    geom        geometry(MultiPolygon, 4326)
);
CREATE INDEX counties_geom_idx ON counties USING GIST (geom);

-- 2 Constituencies (Geometry included)
CREATE TABLE constituencies (
    constituency_code VARCHAR(50) PRIMARY KEY,
    constituency_name VARCHAR(100) NOT NULL,
    county_code VARCHAR(50) NOT NULL REFERENCES counties(county_code) ON DELETE CASCADE,
    geom geometry(MultiPolygon, 4326)
);
CREATE INDEX constituencies_geom_idx ON constituencies USING GIST (geom);
CREATE INDEX constituencies_county_id_idx ON constituencies (county_id);
