import { MapContainer, TileLayer, GeoJSON } from "react-leaflet";
import { useEffect, useState } from "react";
import axios from "axios";
import type { FeatureCollection } from "geojson";
import "leaflet/dist/leaflet.css";

const API_BASE = "http://localhost:8080/api/v1";

export default function KenyaMap() {
  const [counties, setCounties] = useState<FeatureCollection | null>(null);

  useEffect(() => {
    axios
      .get<FeatureCollection>(`${API_BASE}/counties`)
      .then((res) => setCounties(res.data))
      .catch((err) => console.error("Failed to load counties:", err));
  }, []);

  return (
    <MapContainer
      center={[0.0236, 37.9062]}
      zoom={6}
      style={{ height: "100vh", width: "100%" }}
    >
      <TileLayer
        attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>'
        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
      />
      {counties && <GeoJSON data={counties} />}
    </MapContainer>
  );
}
