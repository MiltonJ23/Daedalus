"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { motion } from "framer-motion";
import {
  LayoutDashboard,
  FolderPlus,
  Factory,
  Bell,
  CreditCard,
  ShoppingCart,
  Box,
  BarChart3,
  Shield,
  Tag,
} from "lucide-react";
import clsx from "clsx";

const navItems = [
  { href: "/dashboard", label: "Dashboard", icon: LayoutDashboard },
  { href: "/projects/new", label: "Nouveau Projet", icon: FolderPlus },
  { href: "/procurement", label: "Sourcing", icon: ShoppingCart },
  { href: "/viewer", label: "Visualiseur 3D", icon: Box },
  { href: "/analytics", label: "Analyse coûts", icon: BarChart3 },
  { href: "/notifications", label: "Notifications", icon: Bell },
  { href: "/pricing", label: "Plans & Pricing", icon: Tag },
  { href: "/billing", label: "Facturation", icon: CreditCard },
  { href: "/admin", label: "Administration", icon: Shield },
];

const sidebarVariants = {
  hidden: { x: -280, opacity: 0 },
  visible: {
    x: 0,
    opacity: 1,
    transition: {
      type: "spring",
      stiffness: 300,
      damping: 30,
      staggerChildren: 0.05,
      delayChildren: 0.15,
    },
  },
};

const itemVariants = {
  hidden: { x: -20, opacity: 0 },
  visible: { x: 0, opacity: 1 },
};

export default function Sidebar() {
  const pathname = usePathname();

  return (
    <motion.aside
      variants={sidebarVariants}
      initial="hidden"
      animate="visible"
      className="fixed left-0 top-0 h-screen w-[260px] glass-dark flex flex-col z-50"
    >
      <motion.div
        variants={itemVariants}
        className="px-6 py-7 border-b border-white/10"
      >
        <Link href="/" className="flex items-center gap-3">
          <motion.div
            whileHover={{ rotate: 15, scale: 1.1 }}
            whileTap={{ scale: 0.95 }}
            transition={{ type: "spring", stiffness: 400 }}
            className="relative w-10 h-10 bg-gradient-to-br from-brand-500 to-brand-400 rounded-xl flex items-center justify-center glow-animate"
          >
            <Factory className="w-5 h-5 text-white" />
          </motion.div>
          <div>
            <span className="text-lg font-bold tracking-widest text-white">
              DAEDALUS
            </span>
            <p className="text-[9px] font-medium tracking-[0.2em] text-brand-400 uppercase">
              Industrial Digital Twins
            </p>
          </div>
        </Link>
      </motion.div>

      <nav className="flex-1 px-4 py-6 space-y-1 overflow-y-auto">
        {navItems.map((item) => {
          const isActive =
            pathname === item.href || pathname.startsWith(item.href + "/");
          return (
            <motion.div key={item.href} variants={itemVariants}>
              <Link href={item.href}>
                <motion.div
                  whileHover={{ x: 4 }}
                  whileTap={{ scale: 0.98 }}
                  className={clsx(
                    "relative flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-medium transition-colors",
                    isActive
                      ? "bg-brand-500/15 text-brand-400"
                      : "text-gray-400 hover:bg-white/5 hover:text-white"
                  )}
                >
                  <item.icon className="w-5 h-5" />
                  {item.label}
                  {isActive && (
                    <motion.div
                      layoutId="nav-active"
                      className="absolute left-0 w-1 h-8 bg-gradient-to-b from-brand-500 to-brand-400 rounded-r-full"
                      transition={{
                        type: "spring",
                        stiffness: 500,
                        damping: 30,
                      }}
                    />
                  )}
                </motion.div>
              </Link>
            </motion.div>
          );
        })}
      </nav>

      <motion.div
        variants={itemVariants}
        className="p-4 border-t border-white/10 space-y-3"
      >
        <div className="px-4 py-2 rounded-xl bg-brand-500/10 border border-brand-500/20">
          <p className="text-[10px] text-brand-400 uppercase tracking-wider font-semibold">
            Plan actuel
          </p>
          <p className="text-sm font-bold text-white">FREE</p>
          <p className="text-[10px] text-gray-400 mt-0.5">1/1 projets · 3/3 runs IA</p>
        </div>
        <div className="px-4 py-2 text-xs text-gray-500">v2.0.0 — SRS v2.0</div>
      </motion.div>
    </motion.aside>
  );
}
