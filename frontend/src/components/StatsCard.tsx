"use client";

import { motion } from "framer-motion";
import { useEffect, useState } from "react";
import type { LucideIcon } from "lucide-react";

interface StatsCardProps {
  icon: LucideIcon;
  label: string;
  value: number;
  suffix?: string;
  color: "brand" | "emerald" | "amber" | "rose";
  delay?: number;
}

const colorMap = {
  brand: {
    bg: "bg-brand-50",
    icon: "text-brand-600",
    ring: "ring-brand-100",
    gradient: "from-brand-500 to-brand-600",
  },
  emerald: {
    bg: "bg-emerald-50",
    icon: "text-emerald-600",
    ring: "ring-emerald-100",
    gradient: "from-emerald-500 to-emerald-600",
  },
  amber: {
    bg: "bg-amber-50",
    icon: "text-amber-600",
    ring: "ring-amber-100",
    gradient: "from-amber-500 to-amber-600",
  },
  rose: {
    bg: "bg-rose-50",
    icon: "text-rose-600",
    ring: "ring-rose-100",
    gradient: "from-rose-500 to-rose-600",
  },
};

function useCountUp(target: number, duration: number = 1200) {
  const [count, setCount] = useState(0);

  useEffect(() => {
    if (target === 0) {
      setCount(0);
      return;
    }

    const start = performance.now();
    const step = (now: number) => {
      const elapsed = now - start;
      const progress = Math.min(elapsed / duration, 1);
      const eased = 1 - Math.pow(1 - progress, 3);
      setCount(Math.round(eased * target));
      if (progress < 1) requestAnimationFrame(step);
    };
    requestAnimationFrame(step);
  }, [target, duration]);

  return count;
}

export default function StatsCard({
  icon: Icon,
  label,
  value,
  suffix = "",
  color,
  delay = 0,
}: StatsCardProps) {
  const colors = colorMap[color];
  const animatedValue = useCountUp(value);

  return (
    <motion.div
      initial={{ opacity: 0, y: 24, scale: 0.95 }}
      animate={{ opacity: 1, y: 0, scale: 1 }}
      transition={{ delay, type: "spring", stiffness: 300, damping: 24 }}
      whileHover={{
        y: -4,
        boxShadow: "0 16px 48px -12px rgba(0,0,0,0.12)",
        transition: { duration: 0.2 },
      }}
      className="relative bg-white rounded-2xl p-5 shadow-soft border border-gray-100/60 overflow-hidden group cursor-default"
    >
      {/* Decorative gradient bar */}
      <motion.div
        className={`absolute top-0 left-0 right-0 h-1 bg-gradient-to-r ${colors.gradient}`}
        initial={{ scaleX: 0, originX: 0 }}
        animate={{ scaleX: 1 }}
        transition={{ delay: delay + 0.3, duration: 0.5, ease: "easeOut" }}
      />

      <div className="flex items-center gap-4">
        <motion.div
          whileHover={{ rotate: 8, scale: 1.1 }}
          transition={{ type: "spring", stiffness: 400 }}
          className={`w-12 h-12 ${colors.bg} rounded-xl flex items-center justify-center ring-2 ${colors.ring}`}
        >
          <Icon className={`w-6 h-6 ${colors.icon}`} />
        </motion.div>

        <div>
          <p className="text-xs font-medium text-gray-400 uppercase tracking-wider mb-0.5">
            {label}
          </p>
          <p className="text-2xl font-bold text-gray-900 tabular-nums">
            {animatedValue.toLocaleString("fr-FR")}
            {suffix && (
              <span className="text-sm font-medium text-gray-400 ml-1">
                {suffix}
              </span>
            )}
          </p>
        </div>
      </div>

      {/* Hover shine effect */}
      <motion.div className="absolute inset-0 bg-gradient-to-r from-transparent via-white/60 to-transparent -translate-x-full group-hover:translate-x-full transition-transform duration-700 ease-in-out" />
    </motion.div>
  );
}
