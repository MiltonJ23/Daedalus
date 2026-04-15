"use client";

import { motion, AnimatePresence } from "framer-motion";
import { CheckCircle, XCircle, Info, X } from "lucide-react";
import { useEffect } from "react";

export type ToastType = "success" | "error" | "info";

interface ToastProps {
  id: string;
  message: string;
  type: ToastType;
  onDismiss: (id: string) => void;
  duration?: number;
}

const iconMap = {
  success: CheckCircle,
  error: XCircle,
  info: Info,
};

const styleMap = {
  success: "bg-emerald-50 border-emerald-200 text-emerald-800",
  error: "bg-red-50 border-red-200 text-red-800",
  info: "bg-brand-50 border-brand-200 text-brand-800",
};

const iconColorMap = {
  success: "text-emerald-500",
  error: "text-red-500",
  info: "text-brand-500",
};

export default function Toast({
  id,
  message,
  type,
  onDismiss,
  duration = 4000,
}: ToastProps) {
  const Icon = iconMap[type];

  useEffect(() => {
    const timer = setTimeout(() => onDismiss(id), duration);
    return () => clearTimeout(timer);
  }, [id, onDismiss, duration]);

  return (
    <motion.div
      layout
      initial={{ opacity: 0, x: 80, scale: 0.9 }}
      animate={{ opacity: 1, x: 0, scale: 1 }}
      exit={{ opacity: 0, x: 80, scale: 0.9 }}
      transition={{ type: "spring", stiffness: 400, damping: 25 }}
      className={`relative flex items-center gap-3 px-4 py-3 rounded-xl border shadow-soft ${styleMap[type]}`}
    >
      <Icon className={`w-5 h-5 flex-shrink-0 ${iconColorMap[type]}`} />
      <p className="text-sm font-medium flex-1">{message}</p>
      <button
        onClick={() => onDismiss(id)}
        className="p-0.5 rounded-md hover:bg-black/5 transition"
      >
        <X className="w-4 h-4 opacity-50" />
      </button>

      {/* Auto-dismiss progress bar */}
      <motion.div
        className="absolute bottom-0 left-0 h-0.5 bg-current opacity-20 rounded-b-xl"
        initial={{ width: "100%" }}
        animate={{ width: "0%" }}
        transition={{ duration: duration / 1000, ease: "linear" }}
      />
    </motion.div>
  );
}

interface ToastContainerProps {
  toasts: Array<{ id: string; message: string; type: ToastType }>;
  onDismiss: (id: string) => void;
}

export function ToastContainer({ toasts, onDismiss }: ToastContainerProps) {
  return (
    <div className="fixed bottom-6 right-6 z-[100] flex flex-col gap-2 max-w-sm">
      <AnimatePresence mode="popLayout">
        {toasts.map((toast) => (
          <Toast key={toast.id} {...toast} onDismiss={onDismiss} />
        ))}
      </AnimatePresence>
    </div>
  );
}
