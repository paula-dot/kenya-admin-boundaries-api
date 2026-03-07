import { MapContainer, TileLayer, GeoJSON, useMap } from "react-leaflet";
import { useEffect, useRef } from "react";
import L from "leaflet";
import type { FeatureCollection } from "geojson";
import type { GeoJSON as LeafletGeoJSON } from "leaflet";
import "leaflet/dist/leaflet.css";

type LayerType = "counties" | "constituencies";

interface KenyaMapProps {
  counties: FeatureCollection | null;
  constituencies: FeatureCollection | null;
  activeLayer: LayerType;
  flyToCode: string | null;
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

    // Find the feature matching this county code
    const feature = counties.features.find(
      (f) => f.properties?.code === code
    );
    if (!feature) return;

    // Create a temporary GeoJSON layer to get the bounds
    const layer = L.geoJSON(feature as GeoJSON.Feature);
    const bounds = layer.getBounds();
    if (bounds.isValid()) {
      map.flyToBounds(bounds, { padding: [30, 30], duration: 0.8 });
    }
  }, [code, counties, map]);

  return null;
}

export default function KenyaMap({
  counties,
  constituencies,
  activeLayer,
  flyToCode,
}: KenyaMapProps) {
  const geoJsonRef = useRef<LeafletGeoJSON | null>(null);
  const currentData = activeLayer === "counties" ? counties : constituencies;

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
      {currentData && (
        <GeoJSON
          key={activeLayer}
          ref={geoJsonRef}
          data={currentData}
          style={() => ({
            color: activeLayer === "counties" ? "#00d4ff" : "#ff6b6b",
            weight: activeLayer === "counties" ? 2 : 1,
            fillColor: activeLayer === "counties" ? "#00d4ff" : "#ff6b6b",
            fillOpacity: 0.08,
          })}
          onEachFeature={(feature, layer) => {
            const name = feature.properties?.name || "Unknown";
            const code = feature.properties?.code || "";
            layer.bindTooltip(`${name} (${code})`, {
              sticky: true,
              className: "map-tooltip",
            });
          }}
        />
      )}
      <FlyToCounty counties={counties} code={flyToCode} />
    </MapContainer>
  );
}
