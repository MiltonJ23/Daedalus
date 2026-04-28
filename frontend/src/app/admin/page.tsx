"use client";

import Sidebar from "@/components/Sidebar";
import { motion } from "framer-motion";
import { Users, Activity, Database, Shield } from "lucide-react";

const kpis = [
  { icon: Users, label: "Utilisateurs", value: "1 248" },
  { icon: Activity, label: "Runs IA / jour", value: "342" },
  { icon: Database, label: "Projets actifs", value: "97" },
  { icon: Shield, label: "Incidents 30j", value: "0" },
];

export default function AdminPage() {
  return (
    <div className="flex min-h-screen text-white">
      <Sidebar />
      <main className="flex-1 ml-[260px] p-8">
        <motion.div initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }}>
          <h1 className="text-3xl font-bold mb-2">Administration</h1>
          <p className="text-gray-400">Pilotage opérationnel — réservé aux comptes ADMIN.</p>
        </motion.div>

        <div className="mt-8 grid md:grid-cols-2 lg:grid-cols-4 gap-4">
          {kpis.map((k, i) => (
            <motion.div
              key={k.label}
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: i * 0.05 }}
              className="glass-dark p-6 rounded-2xl"
            >
              <k.icon className="w-5 h-5 text-brand-400 mb-2" />
              <p className="text-xs text-gray-400 uppercase tracking-widest">{k.label}</p>
              <p className="text-2xl font-bold mt-1">{k.value}</p>
            </motion.div>
          ))}
        </div>

        <div className="mt-6 glass-dark p-6 rounded-2xl">
          <h3 className="font-bold mb-4">Liens rapides</h3>
          <ul className="space-y-2 text-sm">
            <li><a href="http://localhost:9090" target="_blank" rel="noopener" className="text-brand-400 hover:underline">→ Prometheus</a></li>
            <li><a href="http://localhost:3001" target="_blank" rel="noopener" className="text-brand-400 hover:underline">→ Grafana</a></li>
            <li><a href="http://localhost:16686" target="_blank" rel="noopener" className="text-brand-400 hover:underline">→ Jaeger</a></li>
            <li><a href="http://localhost:15672" target="_blank" rel="noopener" className="text-brand-400 hover:underline">→ RabbitMQ Management</a></li>
          </ul>
        </div>
      </main>
    </div>
  );
}
