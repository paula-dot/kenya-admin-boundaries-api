import React from "react";

interface DocLayoutProps {
  children: React.ReactNode;
  noPadding?: boolean;
}

export default function DocLayout({ children, noPadding = false }: DocLayoutProps) {
  return (
    <div className="flex h-screen w-screen overflow-hidden bg-background text-foreground font-sans">
      <main className="flex-1 relative overflow-y-auto">
        {noPadding ? (
          <div className="h-full w-full">
            {children}
          </div>
        ) : (
          <div className="max-w-5xl mx-auto px-8 py-12 lg:px-12">
            {children}
          </div>
        )}
      </main>
    </div>
  );
}
