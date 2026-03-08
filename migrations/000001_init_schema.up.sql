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

-- Indexes
-- CREATE INDEX counties_geom_idx ON public.counties USING GIST (geom);
CREATE INDEX constituencies_geom_idx ON constituencies USING GIST (geom);
CREATE INDEX constituencies_county_code_idx ON constituencies (county_code);

CREATE TABLE sub_counties (
    county_code VARCHAR(50) NOT NULL REFERENCES counties(county_code) ON DELETE CASCADE,
    county_name VARCHAR(100) NOT NULL,
    sub_county_code VARCHAR(50) NOT NULL,
    sub_county_name VARCHAR(100) NOT NULL,
    
    -- This tells Postgres: "The code '1' can repeat, but 'KE001'+'1' must be unique"
    PRIMARY KEY (county_code, sub_county_code)
);

CREATE INDEX sub_counties_county_code_idx ON sub_counties (county_code);

-- Wards (without geometry)
CREATE TABLE wards (
    county_code INTEGER NOT NULL,
    county_name VARCHAR(255) NOT NULL,
    constituency_code INTEGER NOT NULL,
    constituency_name VARCHAR(255) NOT NULL,
    ward_code INTEGER PRIMARY KEY,
    ward_name VARCHAR(255) NOT NULL
);
