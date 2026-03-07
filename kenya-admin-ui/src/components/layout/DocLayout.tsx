import React from "react";
import NavSidebar from "../NavSidebar";

export default function DocLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex h-screen w-screen overflow-hidden bg-background text-foreground font-sans">
      <aside className="w-[280px] min-w-[280px] h-full flex-shrink-0 border-r border-border bg-card">
        <NavSidebar />
      </aside>
      <main className="flex-1 relative overflow-y-auto">
        <div className="max-w-5xl mx-auto px-8 py-12 lg:px-12">
          {children}
        </div>
      </main>
    </div>
  );
}
