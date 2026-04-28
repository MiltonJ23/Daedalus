"use client";

import Link from "next/link";
import { motion } from "framer-motion";
import { Factory } from "lucide-react";

export default function RegisterPage() {
  return (
    <main className="min-h-screen flex items-center justify-center px-4 text-white">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="w-full max-w-md glass-dark p-8 rounded-2xl"
      >
        <div className="flex items-center justify-center gap-3 mb-6">
          <div className="w-10 h-10 bg-gradient-to-br from-brand-500 to-brand-400 rounded-xl flex items-center justify-center">
            <Factory className="w-5 h-5 text-white" />
          </div>
          <span className="text-xl font-bold tracking-widest">DAEDALUS</span>
        </div>
        <h1 className="text-2xl font-bold text-center">Créer un compte</h1>
        <p className="mt-2 text-sm text-gray-400 text-center">
          14 jours d&apos;essai sur le plan Business
        </p>

        <form className="mt-8 space-y-4">
          <div>
            <label className="block text-xs text-gray-400 mb-1">Nom complet</label>
            <input type="text" className="w-full px-4 py-3 rounded-xl bg-white/5 border border-white/10 text-white focus:outline-none focus:border-brand-500" placeholder="Jean Dupont" />
          </div>
          <div>
            <label className="block text-xs text-gray-400 mb-1">Email</label>
            <input type="email" className="w-full px-4 py-3 rounded-xl bg-white/5 border border-white/10 text-white focus:outline-none focus:border-brand-500" placeholder="vous@exemple.cm" />
          </div>
          <div>
            <label className="block text-xs text-gray-400 mb-1">Mot de passe</label>
            <input type="password" className="w-full px-4 py-3 rounded-xl bg-white/5 border border-white/10 text-white focus:outline-none focus:border-brand-500" placeholder="Minimum 12 caractères" />
          </div>
          <button type="submit" className="w-full px-4 py-3 rounded-xl bg-brand-500 hover:bg-brand-400 text-white font-semibold">
            Créer mon compte
          </button>
        </form>

        <p className="mt-6 text-sm text-center text-gray-400">
          Déjà inscrit ? <Link href="/login" className="text-brand-400 hover:underline">Se connecter</Link>
        </p>
      </motion.div>
    </main>
  );
}
