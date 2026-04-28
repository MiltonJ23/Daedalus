"use client";

import Sidebar from "@/components/Sidebar";
import { motion } from "framer-motion";
import { ShoppingCart, Package, MapPin } from "lucide-react";

const materials = [
  { name: "Acier IPN 200", qty: "12 tonnes", supplier: "Cameroun Steel — Douala", price: "8 400 000 XAF" },
  { name: "Béton C30/37", qty: "180 m³", supplier: "Cimaf — Yaoundé", price: "16 200 000 XAF" },
  { name: "Toiture bac acier", qty: "1 200 m²", supplier: "Toitures CMR", price: "4 800 000 XAF" },
];

export default function ProcurementPage() {
  return (
    <div className="flex min-h-screen text-white">
      <Sidebar />
      <main className="flex-1 ml-[260px] p-8">
        <motion.div initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }}>
          <h1 className="text-3xl font-bold mb-2">Sourcing automatisé</h1>
          <p className="text-gray-400">Liste de matériaux et fournisseurs locaux trouvés par l&apos;agent procurement.</p>
        </motion.div>

        <div className="mt-8 glass-dark rounded-2xl overflow-hidden">
          <table className="w-full">
            <thead className="bg-white/5 text-left text-xs uppercase tracking-widest text-gray-400">
              <tr>
                <th className="px-6 py-4">Matériau</th>
                <th className="px-6 py-4">Quantité</th>
                <th className="px-6 py-4">Fournisseur</th>
                <th className="px-6 py-4 text-right">Estimation</th>
              </tr>
            </thead>
            <tbody>
              {materials.map((m, i) => (
                <motion.tr key={m.name} initial={{ opacity: 0 }} animate={{ opacity: 1 }} transition={{ delay: i * 0.05 }} className="border-t border-white/5">
                  <td className="px-6 py-4 flex items-center gap-3">
                    <Package className="w-4 h-4 text-brand-400" />{m.name}
                  </td>
                  <td className="px-6 py-4 text-gray-300">{m.qty}</td>
                  <td className="px-6 py-4 text-gray-300 flex items-center gap-2">
                    <MapPin className="w-3 h-3 text-gray-500" /> {m.supplier}
                  </td>
                  <td className="px-6 py-4 text-right font-semibold">{m.price}</td>
                </motion.tr>
              ))}
            </tbody>
          </table>
          <div className="px-6 py-4 border-t border-white/5 flex items-center justify-between">
            <span className="text-sm text-gray-400 flex items-center gap-2">
              <ShoppingCart className="w-4 h-4" /> 3 articles
            </span>
            <span className="font-bold">Total : 29 400 000 XAF</span>
          </div>
        </div>
      </main>
    </div>
  );
}
