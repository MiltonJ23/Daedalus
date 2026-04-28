"use client";

import Link from "next/link";
import { motion } from "framer-motion";
import { Check, Factory } from "lucide-react";

const plans = [
  {
    name: "FREE",
    price: "0",
    tagline: "Pour découvrir Daedalus",
    features: ["1 projet", "3 runs IA / mois", "Export PNG", "Support communautaire"],
    cta: "Commencer",
    highlighted: false,
  },
  {
    name: "STARTER",
    price: "9 900",
    tagline: "Pour les indépendants",
    features: ["5 projets", "30 runs IA / mois", "Export PDF & 3D", "Support email"],
    cta: "Souscrire",
    highlighted: false,
  },
  {
    name: "BUSINESS",
    price: "29 900",
    tagline: "Pour les studios & PME",
    features: [
      "20 projets",
      "150 runs IA / mois",
      "Sourcing automatisé",
      "Analyse de coûts avancée",
      "Support prioritaire",
    ],
    cta: "Souscrire",
    highlighted: true,
  },
  {
    name: "ENTERPRISE",
    price: "Sur devis",
    tagline: "Pour les grandes structures",
    features: [
      "Projets illimités",
      "Runs IA illimités",
      "SSO & audit logs",
      "SLA dédié",
      "Account manager",
    ],
    cta: "Nous contacter",
    highlighted: false,
  },
];

export default function PricingPage() {
  return (
    <main className="min-h-screen text-white">
      <header className="flex items-center justify-between px-8 py-6">
        <Link href="/" className="flex items-center gap-3">
          <div className="w-10 h-10 bg-gradient-to-br from-brand-500 to-brand-400 rounded-xl flex items-center justify-center">
            <Factory className="w-5 h-5 text-white" />
          </div>
          <span className="text-xl font-bold tracking-widest">DAEDALUS</span>
        </Link>
        <Link href="/dashboard" className="text-sm text-gray-300 hover:text-white">
          Dashboard
        </Link>
      </header>

      <section className="max-w-6xl mx-auto px-8 py-12 text-center">
        <h1 className="text-4xl md:text-5xl font-bold">Plans & Tarification</h1>
        <p className="mt-4 text-gray-400">
          Tarifs en XAF, sans engagement. Annulation à tout moment.
        </p>
      </section>

      <section className="max-w-7xl mx-auto px-8 pb-24 grid md:grid-cols-2 lg:grid-cols-4 gap-6">
        {plans.map((p, i) => (
          <motion.div
            key={p.name}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: i * 0.08 }}
            className={`relative glass-dark p-8 rounded-2xl flex flex-col ${
              p.highlighted ? "border-2 border-brand-500" : ""
            }`}
          >
            {p.highlighted && (
              <span className="absolute -top-3 right-6 px-3 py-1 bg-brand-500 text-white text-xs font-bold rounded-full">
                Recommandé
              </span>
            )}
            <h3 className="text-xl font-bold tracking-widest text-brand-400">{p.name}</h3>
            <p className="mt-1 text-xs text-gray-400">{p.tagline}</p>
            <div className="mt-6">
              <span className="text-4xl font-bold">{p.price}</span>
              {p.price !== "Sur devis" && (
                <span className="text-gray-400 text-sm"> XAF / mois</span>
              )}
            </div>
            <ul className="mt-6 space-y-3 flex-1">
              {p.features.map((f) => (
                <li key={f} className="flex items-start gap-2 text-sm text-gray-300">
                  <Check className="w-4 h-4 text-brand-400 mt-0.5 flex-shrink-0" />
                  <span>{f}</span>
                </li>
              ))}
            </ul>
            <Link
              href="/register"
              className={`mt-8 block text-center px-4 py-3 rounded-xl font-semibold ${
                p.highlighted
                  ? "bg-brand-500 hover:bg-brand-400 text-white"
                  : "bg-white/5 hover:bg-white/10 text-white"
              }`}
            >
              {p.cta}
            </Link>
          </motion.div>
        ))}
      </section>
    </main>
  );
}
