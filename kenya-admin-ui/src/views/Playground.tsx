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
    <div className="flex flex-col h-full w-full animate-in fade-in duration-500">
      <header className="mb-6">
        <h1 className="text-3xl font-bold tracking-tight mb-2">Interactive Playground</h1>
        <p className="text-muted-foreground">
          Explore the GeoJSON boundary data returned by our spatial endpoints. Click on a county or constituency to view its properties.
        </p>
      </header>

      <div className="relative flex-1 rounded-xl overflow-hidden border shadow-sm flex flex-col min-h-[500px]">
        <LayerTabs activeLayer={activeLayer} onLayerChange={setActiveLayer} />
        <div className="flex-1 relative z-0">
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
    </div>
  );
}
