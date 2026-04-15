"use client";

import { motion } from "framer-motion";
import { Save, CheckCircle, AlertCircle } from "lucide-react";

interface AutoSaveIndicatorProps {
  lastSavedTime: string | null;
  isSaving: boolean;
  error: string | null;
}

export default function AutoSaveIndicator({
  lastSavedTime,
  isSaving,
  error,
}: AutoSaveIndicatorProps) {
  if (error) {
    return (
      <motion.div
        initial={{ opacity: 0, y: -10 }}
        animate={{ opacity: 1, y: 0 }}
        className="flex items-center gap-2 text-xs text-red-500"
      >
        <AlertCircle className="w-3.5 h-3.5" />
        <span>Erreur de sauvegarde</span>
      </motion.div>
    );
  }

  if (isSaving) {
    return (
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        className="flex items-center gap-2 text-xs text-gray-400 saving-pulse"
      >
        <Save className="w-3.5 h-3.5" />
        <span>Sauvegarde en cours…</span>
      </motion.div>
    );
  }

  if (lastSavedTime) {
    return (
      <motion.div
        initial={{ opacity: 0, y: -10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ type: "spring", stiffness: 500 }}
        className="flex items-center gap-2 text-xs text-brand-600"
      >
        <CheckCircle className="w-3.5 h-3.5" />
        <span>Sauvegardé à {lastSavedTime}</span>
      </motion.div>
    );
  }

  return null;
}
