"use client";

import { useEffect, useState, useCallback } from "react";
import { useParams, useRouter } from "next/navigation";
import { motion } from "framer-motion";
import {
  ArrowLeft,
  MapPin,
  Banknote,
  Ruler,
  Factory,
  Sparkles,
  Save,
  Clock,
} from "lucide-react";
import Sidebar from "@/components/Sidebar";
import AutoSaveIndicator from "@/components/AutoSaveIndicator";
import StatusBadge from "@/components/StatusBadge";
import PageTransition from "@/components/PageTransition";
import { useAutoSave } from "@/hooks/useAutoSave";
import { api } from "@/lib/api";
import type { Project, ProjectUpdate } from "@/types/project";

const INDUSTRIES = [
  "Agroalimentaire",
  "Bois & Menuiserie",
  "Pétrole & Gaz",
  "Mines & Carrières",
  "BTP & Construction",
  "Textile & Confection",
  "Chimie & Cosmétique",
  "Logistique & Transport",
  "Énergie",
  "Métallurgie",
  "Pharmaceutique",
  "Électronique",
  "Autre",
];

const fieldVariants = {
  hidden: { opacity: 0, y: 16 },
  visible: (i: number) => ({
    opacity: 1,
    y: 0,
    transition: {
      delay: 0.15 + i * 0.06,
      type: "spring",
      stiffness: 300,
      damping: 24,
    },
  }),
};

export default function ProjectDetailPage() {
  const params = useParams();
  const router = useRouter();
  const projectId = params.id as string;

  const [project, setProject] = useState<Project | null>(null);
  const [loading, setLoading] = useState(true);
  const [form, setForm] = useState<ProjectUpdate>({});

  const loadProject = useCallback(async () => {
    try {
      const data = await api.projects.get(projectId);
      setProject(data);
      setForm({
        name: data.name,
        industry_type: data.industry_type,
        location: data.location,
        budget: data.budget,
        floor_width: data.floor_width,
        floor_depth: data.floor_depth,
        target_capacity: data.target_capacity,
      });
    } catch {
      router.push("/dashboard");
    } finally {
      setLoading(false);
    }
  }, [projectId, router]);

  useEffect(() => {
    loadProject();
  }, [loadProject]);

  const {
    formattedTime,
    isSaving,
    error: saveError,
    saveNow,
  } = useAutoSave({
    projectId,
    data: form,
    intervalMs: 60_000,
    enabled: !loading && !!project,
  });

  const update = (field: keyof ProjectUpdate, value: string | number) => {
    setForm((prev) => ({ ...prev, [field]: value }));
  };

  const inputClass =
    "w-full px-4 py-3 rounded-xl border border-gray-200 bg-gray-50/50 text-sm focus:outline-none focus:ring-2 focus:ring-brand-300 focus:border-brand-400 transition-all hover:border-gray-300";

  if (loading) {
    return (
      <div className="flex min-h-screen">
        <Sidebar />
        <main className="flex-1 ml-[260px] flex items-center justify-center">
          <div className="flex flex-col items-center gap-4">
            <motion.div
              animate={{ rotate: 360 }}
              transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
              className="w-10 h-10 border-3 border-brand-200 border-t-brand-600 rounded-full"
            />
            <motion.p
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              transition={{ delay: 0.5 }}
              className="text-sm text-gray-400"
            >
              Chargement du projet…
            </motion.p>
          </div>
        </main>
      </div>
    );
  }

  if (!project) return null;

  const area = (form.floor_width ?? 0) * (form.floor_depth ?? 0);

  return (
    <div className="flex min-h-screen">
      <Sidebar />

      <main className="flex-1 ml-[260px]">
        <PageTransition>
          <div className="max-w-3xl mx-auto px-8 py-10">
            {/* Top bar */}
            <div className="flex items-center justify-between mb-8">
              <motion.button
                initial={{ opacity: 0, x: -10 }}
                animate={{ opacity: 1, x: 0 }}
                whileHover={{ x: -4 }}
                onClick={() => router.push("/dashboard")}
                className="flex items-center gap-2 text-sm text-gray-400 hover:text-gray-600 transition"
              >
                <ArrowLeft className="w-4 h-4" />
                Dashboard
              </motion.button>

              <motion.div
                initial={{ opacity: 0, y: -10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.1 }}
                className="flex items-center gap-4"
              >
                <AutoSaveIndicator
                  lastSavedTime={formattedTime}
                  isSaving={isSaving}
                  error={saveError}
                />
                <StatusBadge status={project.status} />
                <motion.button
                  whileHover={{
                    scale: 1.05,
                    boxShadow: "0 4px 15px -2px rgba(0, 212, 170, 0.3)",
                  }}
                  whileTap={{ scale: 0.95 }}
                  onClick={saveNow}
                  className="flex items-center gap-2 px-4 py-2 rounded-xl bg-brand-50 text-brand-700 text-xs font-medium hover:bg-brand-100 transition"
                >
                  <Save className="w-3.5 h-3.5" />
                  Sauvegarder
                </motion.button>
              </motion.div>
            </div>

            {/* Project header info */}
            <motion.div
              initial={{ opacity: 0, y: 16 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.05 }}
              className="mb-6"
            >
              <h1 className="text-2xl font-bold text-gray-900 mb-1">
                {form.name || "Sans nom"}
              </h1>
              <div className="flex items-center gap-3 text-xs text-gray-400">
                <span className="flex items-center gap-1">
                  <Clock className="w-3.5 h-3.5" />
                  Version {project.version}
                </span>
                {area > 0 && (
                  <span className="px-2 py-0.5 bg-brand-50 text-brand-600 rounded-md font-medium">
                    {area.toLocaleString("fr-FR")} m²
                  </span>
                )}
              </div>
            </motion.div>

            {/* Form card */}
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.1 }}
              className="bg-white rounded-2xl shadow-soft p-8 border border-gray-100 space-y-6"
            >
              {/* Name */}
              <motion.div
                custom={0}
                variants={fieldVariants}
                initial="hidden"
                animate="visible"
              >
                <label className="flex items-center gap-2 text-sm font-medium text-gray-700 mb-2">
                  <Factory className="w-4 h-4 text-brand-500" />
                  Nom du projet
                </label>
                <input
                  type="text"
                  value={form.name || ""}
                  onChange={(e) => update("name", e.target.value)}
                  className={inputClass}
                />
              </motion.div>

              {/* Industry */}
              <motion.div
                custom={1}
                variants={fieldVariants}
                initial="hidden"
                animate="visible"
              >
                <label className="flex items-center gap-2 text-sm font-medium text-gray-700 mb-2">
                  <Sparkles className="w-4 h-4 text-brand-500" />
                  Type d&apos;industrie
                </label>
                <select
                  value={form.industry_type || ""}
                  onChange={(e) => update("industry_type", e.target.value)}
                  className={inputClass}
                >
                  <option value="">Sélectionner une industrie</option>
                  {INDUSTRIES.map((opt) => (
                    <option key={opt} value={opt}>
                      {opt}
                    </option>
                  ))}
                </select>
              </motion.div>

              {/* Location */}
              <motion.div
                custom={2}
                variants={fieldVariants}
                initial="hidden"
                animate="visible"
              >
                <label className="flex items-center gap-2 text-sm font-medium text-gray-700 mb-2">
                  <MapPin className="w-4 h-4 text-brand-500" />
                  Localisation
                </label>
                <input
                  type="text"
                  value={form.location || ""}
                  onChange={(e) => update("location", e.target.value)}
                  className={inputClass}
                />
              </motion.div>

              {/* Budget */}
              <motion.div
                custom={3}
                variants={fieldVariants}
                initial="hidden"
                animate="visible"
              >
                <label className="flex items-center gap-2 text-sm font-medium text-gray-700 mb-2">
                  <Banknote className="w-4 h-4 text-brand-500" />
                  Budget (FCFA)
                </label>
                <input
                  type="number"
                  value={form.budget || ""}
                  onChange={(e) =>
                    update("budget", parseFloat(e.target.value) || 0)
                  }
                  className={inputClass}
                />
              </motion.div>

              {/* Dimensions */}
              <motion.div
                custom={4}
                variants={fieldVariants}
                initial="hidden"
                animate="visible"
              >
                <label className="flex items-center gap-2 text-sm font-medium text-gray-700 mb-3">
                  <Ruler className="w-4 h-4 text-brand-500" />
                  Dimensions (m)
                </label>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <input
                      type="number"
                      step="0.1"
                      placeholder="Largeur"
                      value={form.floor_width || ""}
                      onChange={(e) =>
                        update("floor_width", parseFloat(e.target.value) || 0)
                      }
                      className={inputClass}
                    />
                  </div>
                  <div>
                    <input
                      type="number"
                      step="0.1"
                      placeholder="Profondeur"
                      value={form.floor_depth || ""}
                      onChange={(e) =>
                        update("floor_depth", parseFloat(e.target.value) || 0)
                      }
                      className={inputClass}
                    />
                  </div>
                </div>
                {area > 0 && (
                  <motion.p
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    className="text-xs text-brand-600 mt-2"
                  >
                    Surface : {area.toLocaleString("fr-FR")} m²
                  </motion.p>
                )}
              </motion.div>
            </motion.div>

            {/* Auto-save info footer */}
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              transition={{ delay: 0.6 }}
              className="mt-4 text-center"
            >
              <p className="text-xs text-gray-300">
                Sauvegarde automatique toutes les 60 secondes
              </p>
            </motion.div>
          </div>
        </PageTransition>
      </main>
    </div>
  );
}
