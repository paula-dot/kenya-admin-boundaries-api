import React from "react";
import { PrismLight as SyntaxHighlighter } from 'react-syntax-highlighter';
import json from 'react-syntax-highlighter/dist/esm/languages/prism/json';
import docco from 'react-syntax-highlighter/dist/esm/styles/prism/vs';

SyntaxHighlighter.registerLanguage('json', json);

interface EndpointProps {
  method: "GET" | "POST";
  path: string;
  description: string;
  responsePayload: string;
  parameters?: { name: string; type: string; description: string }[];
}

const Endpoint: React.FC<EndpointProps> = ({ method, path, description, responsePayload, parameters }) => {
  const methodColor = method === "GET" ? "text-blue-600 bg-blue-100" : "text-green-600 bg-green-100";
  
  return (
    <div className="mb-16">
      <h3 className="text-xl font-semibold mb-2 flex items-center gap-3">
        <span className={`px-2 py-1 rounded text-xs font-bold font-mono tracking-wider ${methodColor}`}>
          {method}
        </span>
        <code className="bg-muted px-2 py-1 rounded-md text-sm border font-mono">
          {path}
        </code>
      </h3>
      <p className="text-muted-foreground mb-6 leading-relaxed">
        {description}
      </p>

      {parameters && parameters.length > 0 && (
        <div className="mb-6">
          <h4 className="text-sm font-semibold uppercase tracking-wider text-muted-foreground mb-3">Parameters</h4>
          <div className="border rounded-md overflow-hidden bg-card">
            <table className="w-full text-sm text-left">
              <thead className="bg-muted/50 border-b">
                <tr>
                  <th className="px-4 py-3 font-medium">Name</th>
                  <th className="px-4 py-3 font-medium">Type</th>
                  <th className="px-4 py-3 font-medium">Description</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {parameters.map((p, idx) => (
                  <tr key={idx} className="hover:bg-muted/30">
                    <td className="px-4 py-3 font-mono text-xs">{p.name}</td>
                    <td className="px-4 py-3 text-muted-foreground">{p.type}</td>
                    <td className="px-4 py-3 text-muted-foreground">{p.description}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}

      <div>
        <h4 className="text-sm font-semibold uppercase tracking-wider text-muted-foreground mb-3">Example Response</h4>
        <div className="rounded-lg overflow-hidden border bg-[#f8f8f8] shadow-sm">
          <SyntaxHighlighter 
            language="json" 
            style={docco}
            customStyle={{ margin: 0, padding: '1.5rem', fontSize: '0.875rem' }}
          >
            {responsePayload}
          </SyntaxHighlighter>
        </div>
      </div>
    </div>
  );
};

export default function ApiDocs() {
  return (
    <div className="space-y-12 animate-in fade-in slide-in-from-bottom-4 duration-500">
      <header className="border-b pb-8">
        <h1 className="text-4xl font-extrabold tracking-tight mb-4">Kenya Admin Boundaries API</h1>
        <p className="text-xl text-muted-foreground max-w-3xl leading-relaxed">
          A production-grade geospatial REST API mapping Kenya's administrative and electoral boundaries 
          down to the constituency and sub-county level. Designed for sub-millisecond spatial queries.
        </p>
      </header>

      <section>
        <h2 className="text-2xl font-bold border-b pb-2 mb-8">Counties</h2>
        
        <Endpoint 
          method="GET"
          path="/api/v1/counties"
          description="Returns a GeoJSON FeatureCollection of all 47 Kenyan counties with their polygons and metadata."
          responsePayload={`{
  "type": "FeatureCollection",
  "features": [
    {
      "type": "Feature",
      "geometry": {
        "type": "Polygon",
        "coordinates": [[[...]]]
      },
      "properties": {
        "id": 1,
        "code": "KE001",
        "name": "Mombasa"
      }
    }
  ]
}`}
        />

        <Endpoint 
          method="GET"
          path="/api/v1/counties/:slug"
          description="Get a specific county boundary as a single GeoJSON Feature."
          parameters={[
            { name: "slug", type: "string", description: "The county code (e.g., KE047 for Nairobi)." }
          ]}
          responsePayload={`{
  "type": "Feature",
  "geometry": {
    "type": "Polygon",
    "coordinates": [[[...]]]
  },
  "properties": {
    "id": 47,
    "code": "KE047",
    "name": "Nairobi"
  }
}`}
        />
        
        <Endpoint 
          method="GET"
          path="/api/v1/counties/:slug/hierarchy"
          description="A fast, lightweight endpoint returning the County code and name tightly coupled with an array of its Constituencies, completely omitting the heavy PostGIS geometries."
          parameters={[
            { name: "slug", type: "string", description: "The county code (e.g., KE001)." }
          ]}
          responsePayload={`{
  "county_code": "KE001",
  "county_name": "Mombasa",
  "constituencies": [
    {
      "code": "001",
      "name": "Changamwe"
    },
    {
      "code": "002",
      "name": "Jomvu"
    }
  ]
}`}
        />
      </section>

      <section className="pt-8">
        <h2 className="text-2xl font-bold border-b pb-2 mb-8">Sub-Counties</h2>
        <Endpoint 
          method="GET"
          path="/api/v1/sub-counties"
          description="Returns a lightweight JSON array describing all administrative sub-counties."
          responsePayload={`{
  "sub_counties": [
    {
      "county_code": "KE001",
      "county_name": "Mombasa",
      "sub_county_code": "001",
      "sub_county_name": "Changamwe"
    }
  ]
}`}
        />
        <Endpoint 
          method="GET"
          path="/api/v1/counties/:slug/sub-counties"
          description="Returns a lightweight JSON array of sub-counties within a specific county."
          parameters={[
            { name: "slug", type: "string", description: "The county code (e.g., KE047 for Nairobi)." }
          ]}
          responsePayload={`{
  "sub_counties": [
    {
      "county_code": "KE047",
      "county_name": "Nairobi",
      "sub_county_code": "275",
      "sub_county_name": "Starehe"
    }
  ]
}`}
        />
      </section>

      <section className="pt-8">
        <h2 className="text-2xl font-bold border-b pb-2 mb-8">Spatial Intersections</h2>
        
        <Endpoint 
          method="POST"
          path="/api/v1/spatial/intersect"
          description="Submit a Lat/Lng coordinate pair to find exactly which administrative boundaries (County, Constituency) the point falls inside. Highly optimized using PostGIS GIST indexes."
          parameters={[]}
          responsePayload={`{
  "type": "FeatureCollection",
  "features": [
    {
      "type": "Feature",
      "geometry": null,
      "properties": {
        "id": 275,
        "name": "Starehe",
        "type": "constituency"
      }
    },
    {
      "type": "Feature",
      "geometry": { "type": "Polygon", "coordinates": [...] },
      "properties": {
        "id": 47,
        "name": "Nairobi",
        "type": "county"
      }
    }
  ]
}`}
        />
      </section>
    </div>
  );
}
