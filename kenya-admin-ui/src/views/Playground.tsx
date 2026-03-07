import { useEffect, useState } from "react";
import axios from "axios";
import type { FeatureCollection } from "geojson";
import KenyaMap from "../components/KenyaMap";
import LayerTabs from "../components/LayerTabs";
import { Link } from "react-router-dom";
import { ArrowLeft } from "lucide-react";

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
      {/* Absolute Header Overlay for Navigation */}
      <div className="absolute top-4 left-4 z-50">
        <Link
          to="/counties"
          className="flex items-center gap-2 px-4 py-2 bg-card/90 backdrop-blur text-card-foreground shadow-lg rounded-full font-medium hover:bg-card hover:scale-105 transition-all outline-none focus:ring-2 focus:ring-primary"
        >
          <ArrowLeft className="w-5 h-5 text-muted-foreground" />
          Back to API Docs
        </Link>
      </div>

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
