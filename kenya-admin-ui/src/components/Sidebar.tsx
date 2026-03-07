import { useEffect, useState } from "react";
import axios from "axios";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";

const API_BASE = "http://localhost:18080/api/v1";

interface CountyFeature {
  type: string;
  id: string;
  properties: {
    code: string;
    name: string;
  };
  geometry: unknown;
}

interface CountiesResponse {
  type: string;
  features: CountyFeature[];
}

interface SubCounty {
  county_code: string;
  county_name: string;
  sub_county_code: string;
  sub_county_name: string;
}

interface SubCountiesResponse {
  sub_counties: SubCounty[];
}

interface SidebarProps {
  onCountySelect: (code: string, name: string) => void;
}

export default function Sidebar({ onCountySelect }: SidebarProps) {
  const [counties, setCounties] = useState<CountyFeature[]>([]);
  const [subCounties, setSubCounties] = useState<Record<string, SubCounty[]>>({});
  const [loadingCode, setLoadingCode] = useState<string | null>(null);

  useEffect(() => {
    axios
      .get<CountiesResponse>(`${API_BASE}/counties`)
      .then((res) => {
        const sorted = [...res.data.features].sort((a, b) =>
          a.properties.name.localeCompare(b.properties.name)
        );
        setCounties(sorted);
      })
      .catch((err) => console.error("Failed to load counties:", err));
  }, []);

  const fetchSubCountiesForCode = (code: string) => {
    if (!subCounties[code] && loadingCode !== code) {
      setLoadingCode(code);
      axios
        .get<SubCountiesResponse>(`${API_BASE}/counties/${code}/sub-counties`)
        .then((res) => {
          setSubCounties((prev) => ({
            ...prev,
            [code]: res.data.sub_counties || [],
          }));
        })
        .catch((err) => console.error("Failed to load sub-counties:", err))
        .finally(() => setLoadingCode(null));
    }
  };

  const handleTriggerClick = (code: string, name: string) => {
    onCountySelect(code, name);
    fetchSubCountiesForCode(code);
  };

  return (
    <div className="h-full flex flex-col bg-card border-r border-border">
      {/* Header */}
      <div className="p-5 border-b border-border">
        <h1 className="text-lg font-bold tracking-tight text-foreground">
          🇰🇪 Kenya Admin Boundaries
        </h1>
        <p className="text-xs text-muted-foreground mt-1">
          47 Counties · Constituencies · Sub-Counties
        </p>
      </div>

      <Separator />

      {/* County List */}
      <div className="flex-1 overflow-y-auto">
        <Accordion className="px-2">
          {counties.map((county) => (
            <AccordionItem
              key={county.properties.code}
              className="border-b border-border/50"
            >
              <AccordionTrigger
                className="py-3 px-2 text-sm font-medium hover:no-underline hover:bg-accent/50 rounded-md transition-colors"
                onClick={() =>
                  handleTriggerClick(
                    county.properties.code,
                    county.properties.name
                  )
                }
              >
                <div className="flex items-center gap-2">
                  <span>{county.properties.name}</span>
                </div>
              </AccordionTrigger>
              <AccordionContent className="px-2 pb-3">
                <div className="space-y-3">
                  <div className="flex items-center gap-2">
                    <span className="text-xs text-muted-foreground">Code:</span>
                    <Badge variant="secondary" className="text-xs">
                      {county.properties.code}
                    </Badge>
                  </div>

                  <Separator />

                  <div>
                    <p className="text-xs font-semibold text-muted-foreground mb-2">
                      Sub-Counties
                      {subCounties[county.properties.code] && (
                        <span className="ml-1">
                          ({subCounties[county.properties.code].length})
                        </span>
                      )}
                    </p>

                    {loadingCode === county.properties.code ? (
                      <p className="text-xs text-muted-foreground animate-pulse">
                        Loading sub-counties...
                      </p>
                    ) : subCounties[county.properties.code] ? (
                      <ul className="space-y-1">
                        {subCounties[county.properties.code].map((sc) => (
                          <li
                            key={sc.sub_county_code}
                            className="text-xs text-foreground/80 pl-3 py-0.5 border-l-2 border-primary/30 hover:border-primary hover:text-foreground transition-colors"
                          >
                            {sc.sub_county_name}
                          </li>
                        ))}
                      </ul>
                    ) : (
                      <p className="text-xs text-muted-foreground italic">
                        Click to load sub-counties
                      </p>
                    )}
                  </div>
                </div>
              </AccordionContent>
            </AccordionItem>
          ))}
        </Accordion>
      </div>
    </div>
  );
}
