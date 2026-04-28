"use client";

import Sidebar from "@/components/Sidebar";
import { motion } from "framer-motion";
import { CreditCard, Receipt, AlertCircle } from "lucide-react";

export default function BillingPage() {
  return (
    <div className="flex min-h-screen text-white">
      <Sidebar />
      <main className="flex-1 ml-[260px] p-8">
        <motion.div initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }}>
          <h1 className="text-3xl font-bold mb-2">Facturation & Abonnement</h1>
          <p className="text-gray-400">Gérez votre plan, vos factures et votre méthode de paiement.</p>
        </motion.div>

        <div className="mt-8 grid lg:grid-cols-3 gap-6">
          <div className="glass-dark p-6 rounded-2xl lg:col-span-2">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-xs uppercase tracking-widest text-brand-400">Plan actuel</p>
                <h2 className="text-2xl font-bold mt-1">FREE</h2>
                <p className="text-sm text-gray-400">1 projet · 3 runs IA / mois</p>
              </div>
              <a href="/pricing" className="px-4 py-2 rounded-xl bg-brand-500 hover:bg-brand-400 text-white font-medium">
                Changer de plan
              </a>
            </div>
          </div>

          <div className="glass-dark p-6 rounded-2xl">
            <div className="flex items-center gap-3 mb-3">
              <CreditCard className="w-5 h-5 text-brand-400" />
              <h3 className="font-bold">Méthode de paiement</h3>
            </div>
            <p className="text-sm text-gray-400">Aucune méthode enregistrée</p>
            <button className="mt-4 w-full px-4 py-2 rounded-xl bg-white/5 hover:bg-white/10 text-sm">Ajouter une carte</button>
          </div>
        </div>

        <div className="mt-6 glass-dark p-6 rounded-2xl">
          <div className="flex items-center gap-3 mb-4">
            <Receipt className="w-5 h-5 text-brand-400" />
            <h3 className="font-bold">Historique de facturation</h3>
          </div>
          <div className="flex items-center gap-3 px-4 py-6 rounded-xl bg-white/5 text-sm text-gray-400">
            <AlertCircle className="w-4 h-4" />
            Aucune facture pour le moment.
          </div>
        </div>
      </main>
    </div>
  );
}
