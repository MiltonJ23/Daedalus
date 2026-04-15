import type { Metadata } from "next";
import "./globals.css";
import Sidebar from "@/components/Sidebar";

export const metadata: Metadata = {
  title: "Daedalus — Plateforme de conception d'usines industrielles",
  description:
    "Concevez et gérez vos projets d'usines industrielles au Cameroun avec Daedalus — Industrial Digital Twins",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="fr">
      <body className="min-h-screen bg-[var(--bg-primary)] overflow-x-hidden">
        <div className="relative">
          {/* Animated background blobs */}
          <div className="fixed inset-0 overflow-hidden pointer-events-none z-0">
            <div className="blob blob-1" />
            <div className="blob blob-2" />
            <div className="blob blob-3" />
          </div>
          <div className="relative z-10">{children}</div>
        </div>
      </body>
    </html>
  );
}
