import { useEffect, useState } from "react";
import axios from "axios";
import type { FeatureCollection } from "geojson";
import Sidebar from "./components/Sidebar";
import KenyaMap from "./components/KenyaMap";
import LayerTabs from "./components/LayerTabs";
import "./index.css";

const API_BASE = "http://localhost:18080/api/v1";
type LayerType = "counties" | "constituencies";

function App() {
  const [counties, setCounties] = useState<FeatureCollection | null>(null);
  const [constituencies, setConstituencies] = useState<FeatureCollection | null>(null);
  const [activeLayer, setActiveLayer] = useState<LayerType>("counties");
  const [flyToCode, setFlyToCode] = useState<string | null>(null);

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

  const handleCountySelect = (code: string, _name: string) => {
    setFlyToCode(code);
    // Reset after a short delay so re-clicking the same county works
    setTimeout(() => setFlyToCode(null), 1000);
  };

  return (
    <div className="flex h-screen w-screen overflow-hidden bg-background text-foreground">
      {/* Sidebar */}
      <aside className="w-[380px] min-w-[380px] h-full flex-shrink-0">
        <Sidebar onCountySelect={handleCountySelect} />
      </aside>

      {/* Map Area */}
      <main className="flex-1 relative">
        <LayerTabs activeLayer={activeLayer} onLayerChange={setActiveLayer} />
        <KenyaMap
          counties={counties}
          constituencies={constituencies}
          activeLayer={activeLayer}
          flyToCode={flyToCode}
        />
      </main>
    </div>
  );
}

export default App;
