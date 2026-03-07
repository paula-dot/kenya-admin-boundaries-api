import { MapContainer, TileLayer, GeoJSON } from "react-leaflet";
import { useEffect, useState } from "react";
import axios from "axios";
import type { FeatureCollection } from "geojson";
import "leaflet/dist/leaflet.css";

const API_BASE = "http://localhost:18080/api/v1";

type LayerType = "counties" | "constituencies";

export default function KenyaMap() {
  const [counties, setCounties] = useState<FeatureCollection | null>(null);
  const [constituencies, setConstituencies] = useState<FeatureCollection | null>(null);
  const [activeLayer, setActiveLayer] = useState<LayerType>("counties");

  useEffect(() => {
    axios
      .get<FeatureCollection>(`${API_BASE}/counties`)
      .then((res) => setCounties(res.data))
      .catch((err) => console.error("Failed to load counties:", err));

    axios
      .get<FeatureCollection>(`${API_BASE}/constituencies`)
      .then((res) => setConstituencies(res.data))
      .catch((err) => console.error("Failed to load constituencies:", err));
  }, []);

  const currentData = activeLayer === "counties" ? counties : constituencies;

  return (
    <div style={{ position: "relative", height: "100vh", width: "100%" }}>
      {/* Layer toggle control */}
      <div
        style={{
          position: "absolute",
          top: 12,
          right: 12,
          zIndex: 1000,
          background: "rgba(15, 15, 30, 0.9)",
          backdropFilter: "blur(8px)",
          borderRadius: 10,
          padding: "10px 14px",
          display: "flex",
          gap: 8,
          boxShadow: "0 4px 20px rgba(0,0,0,0.4)",
          border: "1px solid rgba(255,255,255,0.08)",
        }}
      >
        <button
          onClick={() => setActiveLayer("counties")}
          style={{
            padding: "6px 16px",
            borderRadius: 6,
            border: "none",
            cursor: "pointer",
            fontWeight: 600,
            fontSize: 13,
            transition: "all 0.2s",
            background: activeLayer === "counties" ? "#00d4ff" : "rgba(255,255,255,0.08)",
            color: activeLayer === "counties" ? "#0a0a1a" : "#ccc",
          }}
        >
          Counties
        </button>
        <button
          onClick={() => setActiveLayer("constituencies")}
          style={{
            padding: "6px 16px",
            borderRadius: 6,
            border: "none",
            cursor: "pointer",
            fontWeight: 600,
            fontSize: 13,
            transition: "all 0.2s",
            background: activeLayer === "constituencies" ? "#00d4ff" : "rgba(255,255,255,0.08)",
            color: activeLayer === "constituencies" ? "#0a0a1a" : "#ccc",
          }}
        >
          Constituencies
        </button>
      </div>

      <MapContainer
        center={[0.0236, 37.9062]}
        zoom={6}
        style={{ height: "100%", width: "100%" }}
      >
        <TileLayer
          attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>'
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        />
        {currentData && (
          <GeoJSON
            key={activeLayer}
            data={currentData}
            style={() => ({
              color: activeLayer === "counties" ? "#00d4ff" : "#ff6b6b",
              weight: activeLayer === "counties" ? 2 : 1,
              fillColor: activeLayer === "counties" ? "#00d4ff" : "#ff6b6b",
              fillOpacity: 0.1,
            })}
          />
        )}
      </MapContainer>
    </div>
  );
}
