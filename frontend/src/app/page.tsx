"use client";

import Link from "next/link";
import { motion } from "framer-motion";
import {
  Factory,
  Sparkles,
  Boxes,
  ShoppingCart,
  BarChart3,
  ArrowRight,
} from "lucide-react";

const features = [
  {
    icon: Sparkles,
    title: "Conception IA",
    desc: "Générez plans 2D et modèles 3D depuis une simple description.",
  },
  {
    icon: Boxes,
    title: "Visualisation 3D",
    desc: "Explorez votre projet en temps réel via React Three Fiber.",
  },
  {
    icon: ShoppingCart,
    title: "Sourcing automatisé",
    desc: "Liste de matériaux et fournisseurs locaux trouvés par l'agent.",
  },
  {
    icon: BarChart3,
    title: "Analyse de coûts",
    desc: "Estimations chiffrées en XAF avec ventilation par poste.",
  },
];

export default function LandingPage() {
  return (
    <main className="min-h-screen text-white">
      <header className="flex items-center justify-between px-8 py-6">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 bg-gradient-to-br from-brand-500 to-brand-400 rounded-xl flex items-center justify-center">
            <Factory className="w-5 h-5 text-white" />
          </div>
          <span className="text-xl font-bold tracking-widest">DAEDALUS</span>
        </div>
        <nav className="flex items-center gap-6 text-sm">
          <Link href="/pricing" className="text-gray-300 hover:text-white">
            Tarifs
          </Link>
          <Link href="/login" className="text-gray-300 hover:text-white">
            Connexion
          </Link>
          <Link
            href="/register"
            className="px-4 py-2 rounded-xl bg-brand-500 hover:bg-brand-400 text-white font-medium"
          >
            Commencer
          </Link>
        </nav>
      </header>

      <section className="max-w-5xl mx-auto px-8 py-24 text-center">
        <motion.h1
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="text-5xl md:text-7xl font-bold leading-tight"
        >
          Concevez vos usines
          <br />
          <span className="bg-gradient-to-r from-brand-400 to-brand-500 bg-clip-text text-transparent">
            avec l&apos;IA
          </span>
        </motion.h1>
        <motion.p
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.15 }}
          className="mt-6 text-lg text-gray-400 max-w-2xl mx-auto"
        >
          Daedalus génère plans 2D, modèles 3D, listes de matériaux et
          estimations de coûts pour vos projets industriels au Cameroun.
        </motion.p>
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.3 }}
          className="mt-10 flex items-center justify-center gap-4"
        >
          <Link
            href="/register"
            className="px-6 py-3 rounded-xl bg-brand-500 hover:bg-brand-400 text-white font-semibold flex items-center gap-2"
          >
            Essai gratuit <ArrowRight className="w-4 h-4" />
          </Link>
          <Link
            href="/pricing"
            className="px-6 py-3 rounded-xl glass-dark text-white font-semibold"
          >
            Voir les plans
          </Link>
        </motion.div>
      </section>

      <section className="max-w-6xl mx-auto px-8 pb-24 grid md:grid-cols-2 lg:grid-cols-4 gap-6">
        {features.map((f, i) => (
          <motion.div
            key={f.title}
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            transition={{ delay: i * 0.08 }}
            viewport={{ once: true }}
            className="glass-dark p-6 rounded-2xl"
          >
            <f.icon className="w-8 h-8 text-brand-400 mb-4" />
            <h3 className="font-bold text-lg mb-2">{f.title}</h3>
            <p className="text-sm text-gray-400">{f.desc}</p>
          </motion.div>
        ))}
      </section>
    </main>
  );
}
