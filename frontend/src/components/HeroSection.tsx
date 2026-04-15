"use client";

import { motion } from "framer-motion";
import { Factory, Cpu, HardHat } from "lucide-react";

const images = [
  {
    src: "https://images.unsplash.com/photo-1565008447742-97f6f38c985c?w=800",
    alt: "Intérieur d'usine industrielle",
    icon: Factory,
  },
  {
    src: "https://images.unsplash.com/photo-1581091226825-a6a2a5aee158?w=800",
    alt: "Technologie industrielle",
    icon: Cpu,
  },
  {
    src: "https://images.unsplash.com/photo-1504307651254-35680f356dfd?w=800",
    alt: "Construction industrielle",
    icon: HardHat,
  },
];

export default function HeroSection() {
  return (
    <motion.section
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      transition={{ duration: 0.6 }}
      className="relative overflow-hidden bg-navy-900 text-white"
    >
      {/* Background gradient overlay */}
      <div className="absolute inset-0 bg-gradient-to-br from-navy-900 via-navy-800 to-brand-900/40" />

      <div className="relative z-10 px-8 py-12">
        <div className="max-w-6xl mx-auto">
          {/* Header text */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.1 }}
            className="mb-8"
          >
            <div className="flex items-center gap-3 mb-3">
              <motion.div
                animate={{ rotate: [0, 5, -5, 0] }}
                transition={{ duration: 4, repeat: Infinity, ease: "easeInOut" }}
                className="w-12 h-12 bg-gradient-to-br from-brand-500 to-brand-400 rounded-xl flex items-center justify-center shadow-glow"
              >
                <Factory className="w-6 h-6 text-white" />
              </motion.div>
              <div>
                <h2 className="text-2xl font-bold tracking-wide">
                  Bienvenue sur{" "}
                  <span className="gradient-text-teal">DAEDALUS</span>
                </h2>
                <p className="text-xs tracking-[0.25em] text-brand-400 font-medium uppercase mt-0.5">
                  Industrial Digital Twins
                </p>
              </div>
            </div>
            <p className="text-sm text-gray-300 max-w-xl leading-relaxed">
              Concevez, modélisez et gérez vos usines industrielles au Cameroun.
              Une plateforme complète pour transformer vos projets en jumeaux
              numériques.
            </p>
          </motion.div>

          {/* Image grid */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            {images.map((img, index) => (
              <motion.div
                key={img.src}
                initial={{ opacity: 0, y: 30 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{
                  delay: 0.2 + index * 0.1,
                  type: "spring",
                  stiffness: 300,
                  damping: 25,
                }}
                whileHover={{ y: -4, scale: 1.02 }}
                className="relative group rounded-xl overflow-hidden cursor-pointer"
              >
                {/* eslint-disable-next-line @next/next/no-img-element */}
                <img
                  src={img.src}
                  alt={img.alt}
                  className="w-full h-40 object-cover transition-transform duration-500 group-hover:scale-110"
                />
                {/* Overlay */}
                <div className="absolute inset-0 bg-gradient-to-t from-navy-900/80 via-navy-900/20 to-transparent" />

                {/* Label */}
                <div className="absolute bottom-3 left-3 flex items-center gap-2">
                  <div className="w-7 h-7 bg-brand-500/20 backdrop-blur-sm rounded-lg flex items-center justify-center border border-brand-400/30">
                    <img.icon className="w-3.5 h-3.5 text-brand-400" />
                  </div>
                  <span className="text-xs font-medium text-white/90">
                    {img.alt}
                  </span>
                </div>

                {/* Hover glow */}
                <motion.div
                  className="absolute inset-0 border-2 border-brand-400/0 rounded-xl group-hover:border-brand-400/40 transition-colors duration-300"
                />
              </motion.div>
            ))}
          </div>
        </div>
      </div>

      {/* Bottom fade */}
      <div className="absolute bottom-0 left-0 right-0 h-8 bg-gradient-to-t from-[var(--bg-primary)] to-transparent" />
    </motion.section>
  );
}
