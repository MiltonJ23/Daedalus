"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { motion } from "framer-motion";
import { ArrowLeft, Sparkles } from "lucide-react";
import Sidebar from "@/components/Sidebar";
import ProjectForm from "@/components/ProjectForm";
import PageTransition from "@/components/PageTransition";
import { api } from "@/lib/api";
import type { ProjectCreate } from "@/types/project";

export default function NewProjectPage() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleCreate = async (data: ProjectCreate) => {
    setLoading(true);
    setError(null);
    try {
      const project = await api.projects.create(data);
      router.push(`/projects/${project.id}`);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Erreur lors de la création");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex min-h-screen">
      <Sidebar />

      <main className="flex-1 ml-[260px]">
        <PageTransition>
          <div className="max-w-2xl mx-auto px-8 py-10">
            {/* Back button */}
            <motion.button
              initial={{ opacity: 0, x: -10 }}
              animate={{ opacity: 1, x: 0 }}
              whileHover={{ x: -4 }}
              onClick={() => router.push("/dashboard")}
              className="flex items-center gap-2 text-sm text-gray-400 hover:text-gray-600 mb-8 transition"
            >
              <ArrowLeft className="w-4 h-4" />
              Retour au dashboard
            </motion.button>

            {/* Title */}
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              className="mb-8"
            >
              <div className="flex items-center gap-3 mb-2">
                <motion.div
                  initial={{ scale: 0 }}
                  animate={{ scale: 1 }}
                  transition={{ delay: 0.2, type: "spring", stiffness: 400 }}
                  className="w-10 h-10 bg-gradient-to-br from-brand-500 to-gold-500 rounded-xl flex items-center justify-center"
                >
                  <Sparkles className="w-5 h-5 text-white" />
                </motion.div>
                <h1 className="text-2xl font-bold text-gray-900">
                  Nouveau projet
                </h1>
              </div>
              <p className="text-sm text-gray-400 ml-[52px]">
                Définissez les paramètres de votre usine pour commencer la
                conception
              </p>
            </motion.div>

            {/* Error */}
            {error && (
              <motion.div
                initial={{ opacity: 0, y: -10 }}
                animate={{ opacity: 1, y: 0 }}
                className="mb-6 p-4 rounded-xl bg-red-50 border border-red-200 text-sm text-red-700"
              >
                {error}
              </motion.div>
            )}

            {/* Card */}
            <motion.div
              initial={{ opacity: 0, y: 30 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.1 }}
              className="bg-white rounded-2xl shadow-soft p-8 border border-gray-100"
            >
              <ProjectForm onSubmit={handleCreate} loading={loading} />
            </motion.div>
          </div>
        </PageTransition>
      </main>
    </div>
  );
}
