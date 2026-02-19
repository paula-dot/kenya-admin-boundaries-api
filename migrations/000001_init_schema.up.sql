-- Enable PostGIS extension
CREATE EXTENSION IF NOT EXISTS postgis;

-- 1. Counties (Geometry included)
CREATE TABLE counties (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    code VARCHAR(10) UNIQUE, -- Official government code
    geom geometry(MultiPolygon, 4326),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
     );
CREATE INDEX counties_geom_idx ON counties USING GIST (geom);

-- 2. Sub-Counties (Relational only, NO geometry)
CREATE TABLE sub_counties (
    id SERIAL PRIMARY KEY,
    county_id INT NOT NULL REFERENCES counties(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    code VARCHAR(10) UNIQUE, -- Official government code
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX sub_counties_county_id_idx ON sub_counties (county_id);

-- 3. Constituencies (Geometry included)
CREATE TABLE constituencies (
    id SERIAL PRIMARY KEY,
    county_id INT NOT NULL REFERENCES counties(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    code VARCHAR(10) UNIQUE, -- Official government code
    geom geometry(MultiPolygon, 4326),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX constituencies_geom_idx ON constituencies USING GIST (geom);
CREATE INDEX constituencies_county_id_idx ON constituencies (county_id);

-- 4. Wards (Geometry included)
CREATE TABLE wards (
    id SERIAL PRIMARY KEY,
    constituency_id INT NOT NULL REFERENCES constituencies(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    code VARCHAR(10) UNIQUE, -- Official government code
    geom geometry(MultiPolygon, 4326),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX wards_geom_idx ON wards USING GIST (geom);
CREATE INDEX wards_constituency_id_idx ON wards (constituency_id);