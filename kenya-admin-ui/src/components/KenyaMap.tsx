import { MapContainer, TileLayer, GeoJSON, useMap } from "react-leaflet";
import { useEffect, useRef, useCallback } from "react";
import L from "leaflet";
import type { FeatureCollection, Feature } from "geojson";
import type { GeoJSON as LeafletGeoJSON, Layer, PathOptions } from "leaflet";
import "leaflet/dist/leaflet.css";

type LayerType = "counties" | "constituencies";

interface KenyaMapProps {
    counties: FeatureCollection | null;
    constituencies: FeatureCollection | null;
    activeLayer: LayerType;
    flyToCode: string | null;
    selectedCode: string | null;
    onFeatureSelect: (code: string, name: string) => void;
}

function FlyToCounty({
    counties,
    code,
}: {
    counties: FeatureCollection | null;
    code: string | null;
}) {
    const map = useMap();

    useEffect(() => {
        if (!code || !counties) return;

        const feature = counties.features.find(
            (f) => f.properties?.code === code
        );
        if (!feature) return;

        const layer = L.geoJSON(feature as GeoJSON.Feature);
        const bounds = layer.getBounds();
        if (bounds.isValid()) {
            map.flyToBounds(bounds, { padding: [30, 30], duration: 0.8 });
        }
    }, [code, counties, map]);

    return null;
}

// Style constants
const COUNTY_COLOR = "#00d4ff";
const CONSTITUENCY_COLOR = "#ff6b6b";

const COUNTY_DEFAULT: PathOptions = {
    color: COUNTY_COLOR,
    weight: 2,
    fillColor: COUNTY_COLOR,
    fillOpacity: 0.08,
};

const COUNTY_HIGHLIGHT: PathOptions = {
    color: COUNTY_COLOR,
    weight: 4,
    fillColor: COUNTY_COLOR,
    fillOpacity: 0.25,
};

// When constituencies are active, county borders stay but are subtler
const COUNTY_BACKDROP: PathOptions = {
    color: COUNTY_COLOR,
    weight: 2,
    fillColor: "transparent",
    fillOpacity: 0,
};

const CONSTITUENCY_DEFAULT: PathOptions = {
    color: CONSTITUENCY_COLOR,
    weight: 1,
    fillColor: CONSTITUENCY_COLOR,
    fillOpacity: 0.08,
};

const CONSTITUENCY_HIGHLIGHT: PathOptions = {
    color: CONSTITUENCY_COLOR,
    weight: 3,
    fillColor: CONSTITUENCY_COLOR,
    fillOpacity: 0.25,
};

export default function KenyaMap({
    counties,
    constituencies,
    activeLayer,
    flyToCode,
    selectedCode,
    onFeatureSelect,
}: KenyaMapProps) {
    const countyRef = useRef<LeafletGeoJSON | null>(null);
    const constituencyRef = useRef<LeafletGeoJSON | null>(null);

    // Update county styles when selectedCode or activeLayer changes
    useEffect(() => {
        if (!countyRef.current) return;

        countyRef.current.eachLayer((layer: Layer) => {
            const feature = (layer as L.GeoJSON & { feature: Feature }).feature;
            const code = feature?.properties?.code;
            const pathLayer = layer as unknown as L.Path;

            if (activeLayer === "counties") {
                pathLayer.setStyle(code === selectedCode ? COUNTY_HIGHLIGHT : COUNTY_DEFAULT);
            } else {
                // Constituencies active: counties are border-only backdrop
                pathLayer.setStyle(COUNTY_BACKDROP);
            }
            if (code === selectedCode) pathLayer.bringToFront();
        });
    }, [selectedCode, activeLayer]);

    // Update constituency styles when selectedCode changes
    useEffect(() => {
        if (!constituencyRef.current || activeLayer !== "constituencies") return;

        constituencyRef.current.eachLayer((layer: Layer) => {
            const feature = (layer as L.GeoJSON & { feature: Feature }).feature;
            const code = feature?.properties?.code;
            const pathLayer = layer as unknown as L.Path;

            pathLayer.setStyle(code === selectedCode ? CONSTITUENCY_HIGHLIGHT : CONSTITUENCY_DEFAULT);
            if (code === selectedCode) pathLayer.bringToFront();
        });
    }, [selectedCode, activeLayer]);

    const handleCountyFeature = useCallback(
        (feature: Feature, layer: Layer) => {
            const name = feature.properties?.name || "Unknown";
            const code = feature.properties?.code || "";

            layer.bindTooltip(`${name} (${code})`, { sticky: true });

            layer.on("click", () => {
                if (activeLayer === "counties") {
                    onFeatureSelect(code, name);
                }
            });

            layer.on("mouseover", () => {
                if (activeLayer === "counties" && code !== selectedCode) {
                    (layer as unknown as L.Path).setStyle({
                        ...COUNTY_DEFAULT,
                        weight: 3,
                        fillOpacity: 0.15,
                    });
                    (layer as unknown as L.Path).bringToFront();
                }
            });

            layer.on("mouseout", () => {
                if (activeLayer === "counties" && code !== selectedCode) {
                    (layer as unknown as L.Path).setStyle(COUNTY_DEFAULT);
                }
            });
        },
        [activeLayer, onFeatureSelect, selectedCode]
    );

    const handleConstituencyFeature = useCallback(
        (feature: Feature, layer: Layer) => {
            const name = feature.properties?.name || "Unknown";
            const code = feature.properties?.code || "";

            layer.bindTooltip(`${name} (${code})`, { sticky: true });

            layer.on("click", () => onFeatureSelect(code, name));

            layer.on("mouseover", () => {
                if (code !== selectedCode) {
                    (layer as unknown as L.Path).setStyle({
                        ...CONSTITUENCY_DEFAULT,
                        weight: 2,
                        fillOpacity: 0.15,
                    });
                    (layer as unknown as L.Path).bringToFront();
                }
            });

            layer.on("mouseout", () => {
                if (code !== selectedCode) {
                    (layer as unknown as L.Path).setStyle(CONSTITUENCY_DEFAULT);
                }
            });
        },
        [onFeatureSelect, selectedCode]
    );

    return (
        <MapContainer
            center={[0.0236, 37.9062]}
            zoom={6}
            style={{ height: "100%", width: "100%" }}
            zoomControl={false}
        >
            <TileLayer
                attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>'
                url="https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}{r}.png"
            />

            {/* Constituency layer (rendered underneath counties when active) */}
            {activeLayer === "constituencies" && constituencies && (
                <GeoJSON
                    key="constituencies"
                    ref={constituencyRef}
                    data={constituencies}
                    style={() => CONSTITUENCY_DEFAULT}
                    onEachFeature={handleConstituencyFeature}
                />
            )}

            {/* County borders — ALWAYS visible */}
            {counties && activeLayer === "counties" && (
                <GeoJSON
                    key="counties-interactive"
                    ref={countyRef}
                    data={counties}
                    style={() => COUNTY_DEFAULT}
                    onEachFeature={handleCountyFeature}
                />
            )}
            {counties && activeLayer === "constituencies" && (
                <GeoJSON
                    key="counties-backdrop"
                    data={counties}
                    style={() => ({ ...COUNTY_BACKDROP, interactive: false })}
                />
            )}

            <FlyToCounty counties={counties} code={flyToCode} />
        </MapContainer>
    );
}
