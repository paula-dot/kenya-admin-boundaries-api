import { useEffect, useState } from "react";
import axios from "axios";
import type { FeatureCollection } from "geojson";
import KenyaMap from "../components/KenyaMap";
import LayerTabs from "../components/LayerTabs";

const API_BASE = "http://localhost:18080/api/v1";
type LayerType = "counties" | "constituencies";

export default function Playground() {
  const [counties, setCounties] = useState<FeatureCollection | null>(null);
  const [constituencies, setConstituencies] = useState<FeatureCollection | null>(null);
  const [activeLayer, setActiveLayer] = useState<LayerType>("counties");
  const [selectedCode, setSelectedCode] = useState<string | null>(null);

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

  const handleFeatureSelect = (code: string, _name: string) => {
    setSelectedCode((prev) => (prev === code ? null : code));
  };

  return (
    <div className="flex flex-col h-full w-full bg-black relative overflow-hidden">
      <LayerTabs activeLayer={activeLayer} onLayerChange={setActiveLayer} />
      <div className="flex-1 relative w-full h-full z-0">
        <KenyaMap
          counties={counties}
          constituencies={constituencies}
          activeLayer={activeLayer}
          flyToCode={null}
          selectedCode={selectedCode}
          onFeatureSelect={handleFeatureSelect}
        />
      </div>
    </div>
  );
}
