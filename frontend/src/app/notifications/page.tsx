"use client";

import Sidebar from "@/components/Sidebar";
import { motion } from "framer-motion";
import { Bell, CheckCheck } from "lucide-react";

const items = [
  { title: "Génération du plan terminée", body: "Le projet « Usine agroalimentaire Douala » a fini sa génération.", time: "il y a 2 min", unread: true },
  { title: "Quota IA bientôt atteint", body: "Plus que 1 run IA disponible ce mois-ci sur le plan FREE.", time: "il y a 1 h", unread: true },
  { title: "Bienvenue sur Daedalus", body: "Découvrez le tutoriel pour créer votre premier projet.", time: "hier", unread: false },
];

export default function NotificationsPage() {
  return (
    <div className="flex min-h-screen text-white">
      <Sidebar />
      <main className="flex-1 ml-[260px] p-8">
        <motion.div initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }} className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold mb-2">Notifications</h1>
            <p className="text-gray-400">Centre de notifications temps réel et historique.</p>
          </div>
          <button className="flex items-center gap-2 px-4 py-2 rounded-xl bg-white/5 hover:bg-white/10 text-sm">
            <CheckCheck className="w-4 h-4" /> Tout marquer comme lu
          </button>
        </motion.div>

        <div className="mt-8 space-y-3">
          {items.map((it, i) => (
            <motion.div
              key={i}
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: i * 0.05 }}
              className={`flex gap-4 p-4 rounded-2xl ${it.unread ? "bg-brand-500/10 border border-brand-500/20" : "glass-dark"}`}
            >
              <div className="w-10 h-10 rounded-xl bg-brand-500/15 flex items-center justify-center flex-shrink-0">
                <Bell className="w-5 h-5 text-brand-400" />
              </div>
              <div className="flex-1">
                <div className="flex items-center justify-between">
                  <h3 className="font-semibold">{it.title}</h3>
                  <span className="text-xs text-gray-500">{it.time}</span>
                </div>
                <p className="text-sm text-gray-400 mt-1">{it.body}</p>
              </div>
            </motion.div>
          ))}
        </div>
      </main>
    </div>
  );
}
