"use client";

import { useState } from "react";
import { motion } from "framer-motion";
import { Factory, MapPin, Banknote, Ruler, Sparkles } from "lucide-react";
import type { ProjectCreate } from "@/types/project";

const industryOptions = [
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

interface ProjectFormProps {
  onSubmit: (data: ProjectCreate) => Promise<void>;
  loading?: boolean;
}

export default function ProjectForm({ onSubmit, loading }: ProjectFormProps) {
  const [form, setForm] = useState<ProjectCreate>({
    name: "",
    industry_type: "",
    location: "",
    budget: 0,
    floor_width: 0,
    floor_depth: 0,
  });
  const [errors, setErrors] = useState<Record<string, string>>({});

  const validate = (): boolean => {
    const errs: Record<string, string> = {};
    if (!form.name.trim()) errs.name = "Le nom est requis";
    if (!form.industry_type)
      errs.industry_type = "Le type d'industrie est requis";
    if (!form.location.trim()) errs.location = "La localisation est requise";
    if (form.budget < 0) errs.budget = "Le budget doit être positif";
    if (!form.floor_width || form.floor_width <= 0)
      errs.floor_width = "La largeur doit être un nombre positif non nul";
    if (!form.floor_depth || form.floor_depth <= 0)
      errs.floor_depth = "La profondeur doit être un nombre positif non nul";
    setErrors(errs);
    return Object.keys(errs).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!validate()) return;
    await onSubmit(form);
  };

  const update = (field: keyof ProjectCreate, value: string | number) => {
    setForm((prev) => ({ ...prev, [field]: value }));
    if (errors[field]) setErrors((prev) => ({ ...prev, [field]: "" }));
  };

  const inputClass = (field: string) =>
    `w-full px-4 py-3 rounded-xl border ${
      errors[field]
        ? "border-red-300 bg-red-50/50"
        : "border-gray-200 bg-gray-50/50"
    } text-sm focus:outline-none focus:ring-2 focus:ring-brand-300 focus:border-brand-400 transition-all`;

  const containerVariants = {
    hidden: { opacity: 0 },
    show: {
      opacity: 1,
      transition: { staggerChildren: 0.1 },
    },
  };

  const itemVariants = {
    hidden: { opacity: 0, y: 20 },
    show: { opacity: 1, y: 0 },
  };

  return (
    <motion.form
      onSubmit={handleSubmit}
      variants={containerVariants}
      initial="hidden"
      animate="show"
      className="space-y-6"
    >
      {/* Project Name */}
      <motion.div variants={itemVariants}>
        <label className="flex items-center gap-2 text-sm font-medium text-gray-700 mb-2">
          <Factory className="w-4 h-4 text-brand-500" />
          Nom du projet
        </label>
        <input
          type="text"
          placeholder="Ex: Usine Cacao Mbalmayo — Ligne de transformation"
          value={form.name}
          onChange={(e) => update("name", e.target.value)}
          className={inputClass("name")}
        />
        {errors.name && (
          <p className="text-xs text-red-500 mt-1">{errors.name}</p>
        )}
      </motion.div>

      {/* Industry type */}
      <motion.div variants={itemVariants}>
        <label className="flex items-center gap-2 text-sm font-medium text-gray-700 mb-2">
          <Sparkles className="w-4 h-4 text-brand-500" />
          Type d&apos;industrie
        </label>
        <select
          value={form.industry_type}
          onChange={(e) => update("industry_type", e.target.value)}
          className={inputClass("industry_type")}
        >
          <option value="">Sélectionner une industrie</option>
          {industryOptions.map((opt) => (
            <option key={opt} value={opt}>
              {opt}
            </option>
          ))}
        </select>
        {errors.industry_type && (
          <p className="text-xs text-red-500 mt-1">{errors.industry_type}</p>
        )}
      </motion.div>

      {/* Location */}
      <motion.div variants={itemVariants}>
        <label className="flex items-center gap-2 text-sm font-medium text-gray-700 mb-2">
          <MapPin className="w-4 h-4 text-brand-500" />
          Localisation
        </label>
        <input
          type="text"
          placeholder="Ex: Douala, Cameroun"
          value={form.location}
          onChange={(e) => update("location", e.target.value)}
          className={inputClass("location")}
        />
        {errors.location && (
          <p className="text-xs text-red-500 mt-1">{errors.location}</p>
        )}
      </motion.div>

      {/* Budget */}
      <motion.div variants={itemVariants}>
        <label className="flex items-center gap-2 text-sm font-medium text-gray-700 mb-2">
          <Banknote className="w-4 h-4 text-brand-500" />
          Budget (FCFA)
        </label>
        <input
          type="number"
          placeholder="50 000 000"
          value={form.budget || ""}
          onChange={(e) => update("budget", parseFloat(e.target.value) || 0)}
          className={inputClass("budget")}
        />
        {errors.budget && (
          <p className="text-xs text-red-500 mt-1">{errors.budget}</p>
        )}
      </motion.div>

      {/* Dimensions */}
      <motion.div variants={itemVariants}>
        <label className="flex items-center gap-2 text-sm font-medium text-gray-700 mb-3">
          <Ruler className="w-4 h-4 text-brand-500" />
          Dimensions de l&apos;usine (m)
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
              className={inputClass("floor_width")}
            />
            {errors.floor_width && (
              <p className="text-xs text-red-500 mt-1">{errors.floor_width}</p>
            )}
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
              className={inputClass("floor_depth")}
            />
            {errors.floor_depth && (
              <p className="text-xs text-red-500 mt-1">{errors.floor_depth}</p>
            )}
          </div>
        </div>
        {form.floor_width > 0 && form.floor_depth > 0 && (
          <motion.p
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className="text-xs text-brand-600 mt-2"
          >
            Surface totale :{" "}
            {(form.floor_width * form.floor_depth).toLocaleString("fr-FR")} m²
          </motion.p>
        )}
      </motion.div>

      {/* Submit */}
      <motion.div variants={itemVariants} className="pt-4">
        <motion.button
          type="submit"
          disabled={loading}
          whileHover={{
            scale: 1.01,
            boxShadow: "0 8px 30px -4px rgba(0, 212, 170, 0.4)",
          }}
          whileTap={{ scale: 0.98 }}
          className="w-full py-3.5 rounded-xl bg-gradient-to-r from-brand-600 to-brand-500 text-white font-semibold text-sm shadow-glow hover:shadow-lg transition-shadow disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {loading ? (
            <span className="flex items-center justify-center gap-2">
              <motion.span
                animate={{ rotate: 360 }}
                transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
                className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full inline-block"
              />
              Création en cours…
            </span>
          ) : (
            "Créer le projet"
          )}
        </motion.button>
      </motion.div>
    </motion.form>
  );
}
