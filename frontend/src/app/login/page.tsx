"use client";

import Link from "next/link";
import { motion } from "framer-motion";
import { Factory, Github, Mail } from "lucide-react";

export default function LoginPage() {
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
        <h1 className="text-2xl font-bold text-center">Connexion</h1>
        <p className="mt-2 text-sm text-gray-400 text-center">
          Accédez à votre espace de conception
        </p>

        <div className="mt-8 space-y-3">
          <button className="w-full flex items-center justify-center gap-3 px-4 py-3 rounded-xl bg-white text-gray-900 font-medium hover:bg-gray-100">
            <svg className="w-5 h-5" viewBox="0 0 24 24">
              <path fill="#4285F4" d="M22.5 12.27c0-.79-.07-1.55-.2-2.27H12v4.3h5.92c-.26 1.38-1.04 2.55-2.21 3.34v2.77h3.57c2.08-1.92 3.27-4.74 3.27-8.14z"/>
              <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.99.66-2.25 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84A10.99 10.99 0 0 0 12 23z"/>
              <path fill="#FBBC04" d="M5.84 14.1A6.6 6.6 0 0 1 5.5 12c0-.73.13-1.44.34-2.1V7.06H2.18A11 11 0 0 0 1 12c0 1.78.43 3.46 1.18 4.94l3.66-2.84z"/>
              <path fill="#EA4335" d="M12 5.38c1.62 0 3.07.56 4.21 1.65l3.15-3.15C17.45 2.1 14.97 1 12 1A10.99 10.99 0 0 0 2.18 7.06l3.66 2.84C6.71 7.31 9.14 5.38 12 5.38z"/>
            </svg>
            Continuer avec Google
          </button>
          <button className="w-full flex items-center justify-center gap-3 px-4 py-3 rounded-xl bg-[#24292e] text-white font-medium hover:bg-[#1c2024]">
            <Github className="w-5 h-5" />
            Continuer avec GitHub
          </button>
        </div>

        <div className="my-6 flex items-center gap-3 text-xs text-gray-500">
          <div className="flex-1 h-px bg-white/10" />ou<div className="flex-1 h-px bg-white/10" />
        </div>

        <form className="space-y-4">
          <div>
            <label className="block text-xs text-gray-400 mb-1">Email</label>
            <input type="email" className="w-full px-4 py-3 rounded-xl bg-white/5 border border-white/10 text-white focus:outline-none focus:border-brand-500" placeholder="vous@exemple.cm" />
          </div>
          <div>
            <label className="block text-xs text-gray-400 mb-1">Mot de passe</label>
            <input type="password" className="w-full px-4 py-3 rounded-xl bg-white/5 border border-white/10 text-white focus:outline-none focus:border-brand-500" placeholder="••••••••" />
          </div>
          <button type="submit" className="w-full flex items-center justify-center gap-2 px-4 py-3 rounded-xl bg-brand-500 hover:bg-brand-400 text-white font-semibold">
            <Mail className="w-4 h-4" /> Se connecter
          </button>
        </form>

        <p className="mt-6 text-sm text-center text-gray-400">
          Pas de compte ? <Link href="/register" className="text-brand-400 hover:underline">Créer un compte</Link>
        </p>
      </motion.div>
    </main>
  );
}
