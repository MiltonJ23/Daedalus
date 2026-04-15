"use client";

import { useState, useCallback } from "react";
import { useRouter } from "next/navigation";
import { motion, AnimatePresence } from "framer-motion";
import {
  Plus,
  Search,
  FolderOpen,
  LayoutGrid,
  Layers,
  TrendingUp,
  Archive,
} from "lucide-react";
import Sidebar from "@/components/Sidebar";
import ProjectCard from "@/components/ProjectCard";
import ConfirmModal from "@/components/ConfirmModal";
import StatsCard from "@/components/StatsCard";
import HeroSection from "@/components/HeroSection";
import PageTransition from "@/components/PageTransition";
import { ToastContainer, type ToastType } from "@/components/Toast";
import { useProjects } from "@/hooks/useProjects";
import { api } from "@/lib/api";

export default function DashboardPage() {
  const router = useRouter();
  const { projects, loading, status, setStatus, refresh } = useProjects();
  const [search, setSearch] = useState("");
  const [deleteTarget, setDeleteTarget] = useState<string | null>(null);
  const [toasts, setToasts] = useState<
    Array<{ id: string; message: string; type: ToastType }>
  >([]);

  const addToast = useCallback((message: string, type: ToastType = "success") => {
    const id = Date.now().toString();
    setToasts((prev) => [...prev, { id, message, type }]);
  }, []);

  const dismissToast = useCallback((id: string) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  const filtered = projects.filter(
    (p) =>
      p.name.toLowerCase().includes(search.toLowerCase()) ||
      p.location.toLowerCase().includes(search.toLowerCase()) ||
      p.industry_type.toLowerCase().includes(search.toLowerCase())
  );

  const activeCount = projects.filter((p) => !p.is_archived).length;
  const archivedCount = projects.filter((p) => p.is_archived).length;
  const totalBudget = projects.reduce(
    (sum, p) => sum + (p.is_archived ? 0 : p.budget),
    0
  );
  const totalArea = projects.reduce(
    (sum, p) => sum + (p.is_archived ? 0 : p.floor_width * p.floor_depth),
    0
  );

  const handleArchive = async (id: string) => {
    await api.projects.archive(id);
    addToast("Projet archivé avec succès", "success");
    refresh();
  };

  const handleRestore = async (id: string) => {
    await api.projects.restore(id);
    addToast("Projet restauré", "info");
    refresh();
  };

  const handleDelete = async () => {
    if (!deleteTarget) return;
    await api.projects.delete(deleteTarget);
    setDeleteTarget(null);
    addToast("Projet supprimé définitivement", "error");
    refresh();
  };

  return (
    <div className="flex min-h-screen">
      <Sidebar />

      <main className="flex-1 ml-[260px]">
        <PageTransition>
          {/* Hero Section */}
          <HeroSection />

          {/* Header */}
          <motion.div
            initial={{ opacity: 0, y: -20 }}
            animate={{ opacity: 1, y: 0 }}
            className="sticky top-0 z-10 glass px-8 py-6 border-b border-gray-100"
          >
            <div className="flex items-center justify-between mb-6">
              <div>
                <motion.h1
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ delay: 0.1 }}
                  className="text-2xl font-bold text-gray-900"
                >
                  Mes projets
                </motion.h1>
                <motion.p
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                  transition={{ delay: 0.2 }}
                  className="text-sm text-gray-400 mt-1"
                >
                  {projects.length} projet{projects.length !== 1 ? "s" : ""} —
                  Tableau de bord
                </motion.p>
              </div>

              <motion.button
                initial={{ opacity: 0, scale: 0.9 }}
                animate={{ opacity: 1, scale: 1 }}
                transition={{ delay: 0.2 }}
                whileHover={{
                  scale: 1.03,
                  boxShadow: "0 8px 30px -4px rgba(0, 212, 170, 0.4)",
                }}
                whileTap={{ scale: 0.97 }}
                onClick={() => router.push("/projects/new")}
                className="flex items-center gap-2 px-5 py-2.5 rounded-xl bg-gradient-to-r from-brand-600 to-brand-500 text-white text-sm font-semibold shadow-glow hover:shadow-lg transition-shadow"
              >
                <Plus className="w-4 h-4" />
                Nouveau projet
              </motion.button>
            </div>

            {/* Filters bar */}
            <div className="flex items-center gap-4">
              {/* Search */}
              <motion.div
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.15 }}
                className="relative flex-1 max-w-md"
              >
                <Search className="absolute left-3.5 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
                <input
                  type="text"
                  placeholder="Rechercher un projet…"
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                  className="w-full pl-10 pr-4 py-2.5 rounded-xl bg-white border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-brand-200 focus:border-brand-400 transition"
                />
              </motion.div>

              {/* Status filter */}
              <motion.div
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.2 }}
                className="flex items-center bg-white rounded-xl border border-gray-200 p-1"
              >
                {[
                  { key: "active", label: "Actifs" },
                  { key: "archived", label: "Archivés" },
                  { key: "all", label: "Tous" },
                ].map((tab) => (
                  <button
                    key={tab.key}
                    onClick={() => setStatus(tab.key)}
                    className={`relative px-4 py-1.5 rounded-lg text-xs font-medium transition-all ${
                      status === tab.key
                        ? "text-brand-700"
                        : "text-gray-500 hover:text-gray-700"
                    }`}
                  >
                    {status === tab.key && (
                      <motion.div
                        layoutId="status-tab"
                        className="absolute inset-0 bg-brand-50 rounded-lg"
                        transition={{
                          type: "spring",
                          stiffness: 500,
                          damping: 30,
                        }}
                      />
                    )}
                    <span className="relative z-10">{tab.label}</span>
                  </button>
                ))}
              </motion.div>
            </div>
          </motion.div>

          {/* Stats Row */}
          {!loading && projects.length > 0 && (
            <div className="px-8 pt-6">
              <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-4">
                <StatsCard
                  icon={LayoutGrid}
                  label="Projets actifs"
                  value={activeCount}
                  color="brand"
                  delay={0.1}
                />
                <StatsCard
                  icon={Archive}
                  label="Archivés"
                  value={archivedCount}
                  color="amber"
                  delay={0.15}
                />
                <StatsCard
                  icon={TrendingUp}
                  label="Budget total"
                  value={Math.round(totalBudget)}
                  suffix=" FCFA"
                  color="emerald"
                  delay={0.2}
                />
                <StatsCard
                  icon={Layers}
                  label="Surface totale"
                  value={Math.round(totalArea)}
                  suffix="m²"
                  color="rose"
                  delay={0.25}
                />
              </div>
            </div>
          )}

          {/* Content */}
          <div className="p-8">
            {loading ? (
              /* Skeleton loaders */
              <motion.div
                initial="hidden"
                animate="show"
                variants={{
                  hidden: {},
                  show: { transition: { staggerChildren: 0.08 } },
                }}
                className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6"
              >
                {[1, 2, 3, 4, 5, 6].map((i) => (
                  <motion.div
                    key={i}
                    variants={{
                      hidden: { opacity: 0, y: 20 },
                      show: { opacity: 1, y: 0 },
                    }}
                    className="bg-white rounded-2xl p-6 shadow-soft"
                  >
                    <div className="shimmer h-5 w-3/4 rounded-lg mb-3" />
                    <div className="shimmer h-3 w-1/2 rounded-lg mb-4" />
                    <div className="flex gap-2">
                      <div className="shimmer h-7 w-24 rounded-lg" />
                      <div className="shimmer h-7 w-20 rounded-lg" />
                      <div className="shimmer h-7 w-16 rounded-lg" />
                    </div>
                  </motion.div>
                ))}
              </motion.div>
            ) : filtered.length === 0 ? (
              /* Empty state */
              <motion.div
                initial={{ opacity: 0, scale: 0.95 }}
                animate={{ opacity: 1, scale: 1 }}
                className="flex flex-col items-center justify-center py-24 text-center"
              >
                <motion.div
                  animate={{ y: [0, -8, 0] }}
                  transition={{
                    duration: 3,
                    repeat: Infinity,
                    ease: "easeInOut",
                  }}
                  className="w-20 h-20 bg-brand-50 rounded-2xl flex items-center justify-center mb-6 ring-4 ring-brand-50/50"
                >
                  <FolderOpen className="w-10 h-10 text-brand-400" />
                </motion.div>
                <h3 className="text-lg font-semibold text-gray-800 mb-2">
                  {search ? "Aucun résultat" : "Aucun projet"}
                </h3>
                <p className="text-sm text-gray-400 mb-6 max-w-sm">
                  {search
                    ? "Essayez un autre terme de recherche"
                    : "Créez votre premier projet pour commencer à concevoir votre usine"}
                </p>
                {!search && (
                  <motion.button
                    whileHover={{
                      scale: 1.03,
                      boxShadow: "0 8px 30px -4px rgba(0, 212, 170, 0.4)",
                    }}
                    whileTap={{ scale: 0.97 }}
                    onClick={() => router.push("/projects/new")}
                    className="flex items-center gap-2 px-5 py-2.5 rounded-xl bg-gradient-to-r from-brand-600 to-brand-500 text-white text-sm font-semibold shadow-glow"
                  >
                    <Plus className="w-4 h-4" />
                    Créer un projet
                  </motion.button>
                )}
              </motion.div>
            ) : (
              /* Project grid */
              <motion.div
                layout
                className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6"
              >
                <AnimatePresence mode="popLayout">
                  {filtered.map((project, index) => (
                    <ProjectCard
                      key={project.id}
                      project={project}
                      index={index}
                      onArchive={handleArchive}
                      onRestore={handleRestore}
                      onDelete={(id) => setDeleteTarget(id)}
                      onClick={(id) => router.push(`/projects/${id}`)}
                    />
                  ))}
                </AnimatePresence>
              </motion.div>
            )}
          </div>
        </PageTransition>
      </main>

      {/* Delete confirmation */}
      <ConfirmModal
        isOpen={!!deleteTarget}
        title="Supprimer le projet"
        message="Cette action est irréversible. Le projet et toutes ses données seront définitivement supprimés."
        confirmLabel="Supprimer définitivement"
        onConfirm={handleDelete}
        onCancel={() => setDeleteTarget(null)}
      />

      {/* Toast notifications */}
      <ToastContainer toasts={toasts} onDismiss={dismissToast} />
    </div>
  );
}
