"use client";

import Sidebar from "@/components/Sidebar";
import { motion } from "framer-motion";
import { TrendingUp, TrendingDown, DollarSign, PieChart } from "lucide-react";

const breakdown = [
  { label: "Gros œuvre", value: 38, amount: "11 200 000 XAF" },
  { label: "Charpente métallique", value: 24, amount: "7 050 000 XAF" },
  { label: "Toiture & étanchéité", value: 16, amount: "4 700 000 XAF" },
  { label: "Électricité", value: 12, amount: "3 530 000 XAF" },
  { label: "Plomberie", value: 10, amount: "2 940 000 XAF" },
];

export default function AnalyticsPage() {
  return (
    <div className="flex min-h-screen text-white">
      <Sidebar />
      <main className="flex-1 ml-[260px] p-8">
        <motion.div initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }}>
          <h1 className="text-3xl font-bold mb-2">Analyse des coûts</h1>
          <p className="text-gray-400">Estimations chiffrées en XAF par poste de dépense.</p>
        </motion.div>

        <div className="mt-8 grid md:grid-cols-3 gap-4">
          <div className="glass-dark p-6 rounded-2xl">
            <DollarSign className="w-5 h-5 text-brand-400 mb-2" />
            <p className="text-xs text-gray-400 uppercase tracking-widest">Total estimé</p>
            <p className="text-2xl font-bold mt-1">29 420 000 XAF</p>
          </div>
          <div className="glass-dark p-6 rounded-2xl">
            <TrendingDown className="w-5 h-5 text-emerald-400 mb-2" />
            <p className="text-xs text-gray-400 uppercase tracking-widest">vs marché</p>
            <p className="text-2xl font-bold mt-1">-12%</p>
          </div>
          <div className="glass-dark p-6 rounded-2xl">
            <TrendingUp className="w-5 h-5 text-amber-400 mb-2" />
            <p className="text-xs text-gray-400 uppercase tracking-widest">Marge sécurité</p>
            <p className="text-2xl font-bold mt-1">+8%</p>
          </div>
        </div>

        <div className="mt-6 glass-dark p-6 rounded-2xl">
          <div className="flex items-center gap-3 mb-6">
            <PieChart className="w-5 h-5 text-brand-400" />
            <h3 className="font-bold">Ventilation par poste</h3>
          </div>
          <div className="space-y-4">
            {breakdown.map((b) => (
              <div key={b.label}>
                <div className="flex items-center justify-between text-sm mb-1">
                  <span>{b.label}</span>
                  <span className="text-gray-400">{b.amount} · {b.value}%</span>
                </div>
                <div className="h-2 rounded-full bg-white/5 overflow-hidden">
                  <motion.div
                    initial={{ width: 0 }}
                    animate={{ width: `${b.value}%` }}
                    transition={{ duration: 0.6 }}
                    className="h-full bg-gradient-to-r from-brand-500 to-brand-400"
                  />
                </div>
              </div>
            ))}
          </div>
        </div>
      </main>
    </div>
  );
}
