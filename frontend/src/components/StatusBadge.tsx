"use client";

import { motion } from "framer-motion";
import clsx from "clsx";

interface StatusBadgeProps {
  status: "active" | "archived" | string;
}

const statusConfig: Record<
  string,
  { label: string; color: string; dot: string }
> = {
  active: {
    label: "Actif",
    color: "bg-emerald-50 text-emerald-700",
    dot: "bg-emerald-500",
  },
  archived: {
    label: "Archivé",
    color: "bg-amber-50 text-amber-700",
    dot: "bg-amber-500",
  },
  draft: {
    label: "Brouillon",
    color: "bg-gray-50 text-gray-600",
    dot: "bg-gray-400",
  },
};

export default function StatusBadge({ status }: StatusBadgeProps) {
  const config = statusConfig[status] || statusConfig.draft;

  return (
    <motion.span
      initial={{ scale: 0.8, opacity: 0 }}
      animate={{ scale: 1, opacity: 1 }}
      className={clsx(
        "inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium",
        config.color
      )}
    >
      <span className={clsx("w-1.5 h-1.5 rounded-full", config.dot)} />
      {config.label}
    </motion.span>
  );
}
