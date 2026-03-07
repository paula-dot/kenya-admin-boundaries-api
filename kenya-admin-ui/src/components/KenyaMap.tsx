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

function getDefaultStyle(activeLayer: LayerType): PathOptions {
  const color = activeLayer === "counties" ? COUNTY_COLOR : CONSTITUENCY_COLOR;
  return {
    color,
    weight: activeLayer === "counties" ? 2 : 1,
    fillColor: color,
    fillOpacity: 0.08,
  };
}

function getHighlightStyle(activeLayer: LayerType): PathOptions {
  const color = activeLayer === "counties" ? COUNTY_COLOR : CONSTITUENCY_COLOR;
  return {
    color,
    weight: 4,
    fillColor: color,
    fillOpacity: 0.3,
  };
}

export default function KenyaMap({
  counties,
  constituencies,
  activeLayer,
  flyToCode,
  selectedCode,
  onFeatureSelect,
}: KenyaMapProps) {
  const geoJsonRef = useRef<LeafletGeoJSON | null>(null);
  const currentData = activeLayer === "counties" ? counties : constituencies;

  // Update styles when selectedCode changes
  useEffect(() => {
    if (!geoJsonRef.current) return;

    geoJsonRef.current.eachLayer((layer: Layer) => {
      const feature = (layer as L.GeoJSON & { feature: Feature }).feature;
      const code = feature?.properties?.code;
      const pathLayer = layer as unknown as L.Path;

      if (code === selectedCode) {
        pathLayer.setStyle(getHighlightStyle(activeLayer));
        pathLayer.bringToFront();
      } else {
        pathLayer.setStyle(getDefaultStyle(activeLayer));
      }
    });
  }, [selectedCode, activeLayer]);

  const handleEachFeature = useCallback(
    (feature: Feature, layer: Layer) => {
      const name = feature.properties?.name || "Unknown";
      const code = feature.properties?.code || "";

      // Tooltip
      layer.bindTooltip(`${name} (${code})`, {
        sticky: true,
        className: "map-tooltip",
      });

      // Click → highlight + notify parent
      layer.on("click", () => {
        onFeatureSelect(code, name);
      });

      // Hover effect (subtle brightening)
      layer.on("mouseover", () => {
        const pathLayer = layer as unknown as L.Path;
        if (code !== selectedCode) {
          pathLayer.setStyle({
            ...getDefaultStyle(activeLayer),
            weight: activeLayer === "counties" ? 3 : 2,
            fillOpacity: 0.15,
          });
          pathLayer.bringToFront();
        }
      });

      layer.on("mouseout", () => {
        const pathLayer = layer as unknown as L.Path;
        if (code !== selectedCode) {
          pathLayer.setStyle(getDefaultStyle(activeLayer));
        }
      });
    },
    [activeLayer, onFeatureSelect, selectedCode]
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
      {currentData && (
        <GeoJSON
          key={activeLayer}
          ref={geoJsonRef}
          data={currentData}
          style={() => getDefaultStyle(activeLayer)}
          onEachFeature={handleEachFeature}
        />
      )}
      <FlyToCounty counties={counties} code={flyToCode} />
    </MapContainer>
  );
}
