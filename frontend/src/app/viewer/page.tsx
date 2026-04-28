"use client";

import Sidebar from "@/components/Sidebar";
import { motion } from "framer-motion";
import { Box, Download, RotateCcw } from "lucide-react";

export default function ViewerPage() {
  return (
    <div className="flex min-h-screen text-white">
      <Sidebar />
      <main className="flex-1 ml-[260px] p-8">
        <motion.div initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }} className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold mb-2">Visualiseur 3D</h1>
            <p className="text-gray-400">Inspectez votre projet (R3F · GLB / GLTF). Données issues de Meshy.ai.</p>
          </div>
          <div className="flex gap-2">
            <button className="flex items-center gap-2 px-4 py-2 rounded-xl bg-white/5 hover:bg-white/10 text-sm">
              <RotateCcw className="w-4 h-4" /> Réinitialiser
            </button>
            <button className="flex items-center gap-2 px-4 py-2 rounded-xl bg-brand-500 hover:bg-brand-400 text-sm font-medium">
              <Download className="w-4 h-4" /> Exporter GLB
            </button>
          </div>
        </motion.div>

        <div className="mt-8 glass-dark rounded-2xl aspect-video flex items-center justify-center text-center p-8">
          <div>
            <Box className="w-16 h-16 text-brand-400 mx-auto mb-4" />
            <h2 className="text-xl font-semibold">Canvas R3F</h2>
            <p className="mt-2 text-sm text-gray-400 max-w-md mx-auto">
              Le canvas Three.js sera monté ici. Sélectionnez un projet depuis le dashboard pour charger son modèle 3D.
            </p>
          </div>
        </div>
      </main>
    </div>
  );
}
