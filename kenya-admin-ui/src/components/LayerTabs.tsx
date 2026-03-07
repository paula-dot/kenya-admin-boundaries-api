import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";

type LayerType = "counties" | "constituencies";

interface LayerTabsProps {
  activeLayer: LayerType;
  onLayerChange: (layer: LayerType) => void;
}

export default function LayerTabs({ activeLayer, onLayerChange }: LayerTabsProps) {
  return (
    <div className="absolute top-3 left-1/2 -translate-x-1/2 z-[1000]">
      <Tabs
        value={activeLayer}
        onValueChange={(v) => onLayerChange(v as LayerType)}
      >
        <TabsList className="bg-card/90 backdrop-blur-md border border-border shadow-lg">
          <TabsTrigger value="counties" className="text-xs font-semibold">
            Counties
          </TabsTrigger>
          <TabsTrigger value="constituencies" className="text-xs font-semibold">
            Constituencies
          </TabsTrigger>
          <TabsTrigger value="sub-counties" className="text-xs font-semibold">
            Sub-Counties
          </TabsTrigger>
        </TabsList>
      </Tabs>
    </div>
  );
}
