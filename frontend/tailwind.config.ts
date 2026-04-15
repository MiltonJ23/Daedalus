import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./src/app/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        brand: {
          50: "#edfffe",
          100: "#d1fefa",
          200: "#a9fbf5",
          300: "#6ef5eb",
          400: "#4ecdc4",
          500: "#00d4aa",
          600: "#00ab8e",
          700: "#008872",
          800: "#066b5d",
          900: "#0a584d",
        },
        navy: {
          50: "#e8ecf2",
          100: "#c5cee0",
          200: "#9eaecc",
          300: "#778eb7",
          400: "#5a76a7",
          500: "#3d5e97",
          600: "#2b4575",
          700: "#1c3054",
          800: "#0f1f38",
          900: "#0a1628",
        },
        gold: {
          50: "#fffbeb",
          100: "#fef3c7",
          200: "#fde68a",
          300: "#fcd34d",
          400: "#fbbf24",
          500: "#f59e0b",
          600: "#d97706",
          700: "#b45309",
          800: "#92400e",
          900: "#78350f",
        },
      },
      fontFamily: {
        sans: ["Inter", "system-ui", "sans-serif"],
      },
      boxShadow: {
        soft: "0 2px 15px -3px rgba(0, 0, 0, 0.07), 0 10px 20px -2px rgba(0, 0, 0, 0.04)",
        glow: "0 0 20px rgba(0, 212, 170, 0.25)",
        "glow-gold": "0 0 20px rgba(245, 158, 11, 0.15)",
      },
    },
  },
  plugins: [],
};

export default config;
