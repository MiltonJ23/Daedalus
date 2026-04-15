"use client";

import { motion } from "framer-motion";
import { formatDistanceToNow } from "date-fns";
import { fr } from "date-fns/locale";
import {
  Ruler,
  MapPin,
  Banknote,
  MoreVertical,
  Archive,
  Trash2,
  RotateCcw,
  ExternalLink,
} from "lucide-react";
import { useState } from "react";
import type { Project } from "@/types/project";
import StatusBadge from "./StatusBadge";

interface ProjectCardProps {
  project: Project;
  index: number;
  onArchive: (id: string) => void;
  onRestore: (id: string) => void;
  onDelete: (id: string) => void;
  onClick: (id: string) => void;
}

export default function ProjectCard({
  project,
  index,
  onArchive,
  onRestore,
  onDelete,
  onClick,
}: ProjectCardProps) {
  const [menuOpen, setMenuOpen] = useState(false);

  const updatedAgo = formatDistanceToNow(new Date(project.updated_at), {
    addSuffix: true,
    locale: fr,
  });

  const area = (project.floor_width * project.floor_depth).toLocaleString(
    "fr-FR"
  );

  return (
    <motion.div
      initial={{ opacity: 0, y: 30 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, y: -20, scale: 0.95 }}
      transition={{
        delay: index * 0.08,
        type: "spring",
        stiffness: 300,
        damping: 25,
      }}
      whileHover={{
        y: -6,
        boxShadow: "0 20px 50px -12px rgba(0,0,0,0.15)",
      }}
      className="relative bg-white rounded-2xl p-6 shadow-soft cursor-pointer group border border-transparent hover:border-brand-100/80 transition-colors overflow-hidden"
      onClick={() => onClick(project.id)}
    >
      {/* Gradient accent stripe */}
      <div className="absolute top-0 left-0 right-0 h-1 bg-gradient-to-r from-brand-400 via-gold-400 to-brand-500 opacity-0 group-hover:opacity-100 transition-opacity duration-300" />

      {/* Header */}
      <div className="flex items-start justify-between mb-4">
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 mb-1">
            <h3 className="text-base font-semibold text-gray-900 truncate">
              {project.name}
            </h3>
            <motion.div
              initial={{ opacity: 0, scale: 0 }}
              whileHover={{ opacity: 1, scale: 1 }}
              className="text-brand-400"
            >
              <ExternalLink className="w-3.5 h-3.5" />
            </motion.div>
          </div>
          <p className="text-xs text-gray-400">Mis à jour {updatedAgo}</p>
        </div>

        <div className="flex items-center gap-2">
          <StatusBadge status={project.status} />

          {/* Menu */}
          <div className="relative">
            <motion.button
              whileHover={{ scale: 1.1 }}
              whileTap={{ scale: 0.9 }}
              onClick={(e) => {
                e.stopPropagation();
                setMenuOpen(!menuOpen);
              }}
              className="p-1.5 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 opacity-0 group-hover:opacity-100 transition-all"
            >
              <MoreVertical className="w-4 h-4" />
            </motion.button>

            {menuOpen && (
              <motion.div
                initial={{ opacity: 0, scale: 0.9, y: -4 }}
                animate={{ opacity: 1, scale: 1, y: 0 }}
                exit={{ opacity: 0, scale: 0.9 }}
                className="absolute right-0 top-8 w-44 bg-white rounded-xl shadow-lg border border-gray-100 py-1 z-10"
                onClick={(e) => e.stopPropagation()}
              >
                {project.is_archived ? (
                  <button
                    onClick={() => {
                      onRestore(project.id);
                      setMenuOpen(false);
                    }}
                    className="flex items-center gap-2 w-full px-4 py-2.5 text-sm text-gray-700 hover:bg-gray-50 transition"
                  >
                    <RotateCcw className="w-4 h-4" /> Restaurer
                  </button>
                ) : (
                  <button
                    onClick={() => {
                      onArchive(project.id);
                      setMenuOpen(false);
                    }}
                    className="flex items-center gap-2 w-full px-4 py-2.5 text-sm text-gray-700 hover:bg-gray-50 transition"
                  >
                    <Archive className="w-4 h-4" /> Archiver
                  </button>
                )}
                <button
                  onClick={() => {
                    onDelete(project.id);
                    setMenuOpen(false);
                  }}
                  className="flex items-center gap-2 w-full px-4 py-2.5 text-sm text-red-600 hover:bg-red-50 transition"
                >
                  <Trash2 className="w-4 h-4" /> Supprimer
                </button>
              </motion.div>
            )}
          </div>
        </div>
      </div>

      {/* Industry type pill */}
      <p className="text-xs text-brand-600 font-medium mb-3">
        {project.industry_type}
      </p>

      {/* Info pills */}
      <div className="flex flex-wrap gap-2">
        <motion.span
          whileHover={{ scale: 1.05 }}
          className="inline-flex items-center gap-1.5 text-xs text-gray-500 bg-gray-50 px-3 py-1.5 rounded-lg hover:bg-gray-100 transition-colors"
        >
          <MapPin className="w-3.5 h-3.5" />
          {project.location}
        </motion.span>
        <motion.span
          whileHover={{ scale: 1.05 }}
          className="inline-flex items-center gap-1.5 text-xs text-gray-500 bg-gray-50 px-3 py-1.5 rounded-lg hover:bg-gray-100 transition-colors"
        >
          <Banknote className="w-3.5 h-3.5" />
          {project.budget.toLocaleString("fr-FR")} FCFA
        </motion.span>
        <motion.span
          whileHover={{ scale: 1.05 }}
          className="inline-flex items-center gap-1.5 text-xs text-gray-500 bg-gray-50 px-3 py-1.5 rounded-lg hover:bg-gray-100 transition-colors"
        >
          <Ruler className="w-3.5 h-3.5" />
          {area} m²
        </motion.span>
      </div>

      {/* Bottom gradient line */}
      <motion.div
        className="absolute bottom-0 left-6 right-6 h-0.5 bg-gradient-to-r from-brand-400 to-gold-400 rounded-full origin-left"
        initial={{ scaleX: 0 }}
        whileHover={{ scaleX: 1 }}
        transition={{ duration: 0.3 }}
      />

      {/* Hover shine effect */}
      <div className="absolute inset-0 bg-gradient-to-r from-transparent via-white/40 to-transparent -translate-x-full group-hover:translate-x-full transition-transform duration-700 ease-in-out pointer-events-none" />
    </motion.div>
  );
}
