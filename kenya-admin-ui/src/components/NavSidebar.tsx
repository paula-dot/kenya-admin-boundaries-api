import { Link, useLocation } from "react-router-dom";
import { Book, Map, Code, Layers, FileJson } from "lucide-react";

export default function NavSidebar() {
  const location = useLocation();

  const links = [
    { name: "Getting Started", path: "/", icon: <Book className="w-4 h-4" /> },
    { name: "Counties Endpoint", path: "/counties", icon: <Layers className="w-4 h-4" /> },
    { name: "Constituencies Endpoint", path: "/constituencies", icon: <FileJson className="w-4 h-4" /> },
    { name: "Sub-Counties Endpoint", path: "/sub-counties", icon: <FileJson className="w-4 h-4" /> },
    { name: "Spatial Intersect", path: "/spatial", icon: <Code className="w-4 h-4" /> },
    { name: "Interactive Map", path: "/map", icon: <Map className="w-4 h-4" /> },
  ];

  return (
    <div className="h-full flex flex-col bg-card border-r border-border">
      <div className="p-5 border-b border-border text-center sm:text-left">
        <h1 className="text-lg font-bold tracking-tight text-foreground">
          🇰🇪 Kenya Boundary API
        </h1>
        <p className="text-xs text-muted-foreground mt-1">Documentation Portal</p>
      </div>

      <nav className="flex-1 overflow-y-auto py-4">
        <ul className="space-y-1 px-4">
          {links.map((link) => {
            const isActive = location.pathname === link.path;
            return (
              <li key={link.name}>
                <Link
                  to={link.path}
                  className={`flex items-center gap-3 px-3 py-2.5 rounded-md transition-colors text-sm font-medium ${
                    isActive
                      ? "bg-primary text-primary-foreground shadow-sm"
                      : "text-muted-foreground hover:bg-accent hover:text-accent-foreground"
                  }`}
                >
                  {link.icon}
                  {link.name}
                </Link>
              </li>
            );
          })}
        </ul>
      </nav>
    </div>
  );
}
