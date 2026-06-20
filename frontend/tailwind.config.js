/** @type {import('tailwindcss').Config} */

export default {
  darkMode: "class",
  content: ["./index.html", "./src/**/*.{js,ts,vue}"],
  theme: {
    extend: {
      colors: {
        'osc-bg': '#0a0e14',
        'osc-panel': '#0f1419',
        'osc-card': '#161b22',
        'osc-border': '#21262d',
        'osc-muted': '#6e7681',
        'osc-text': '#c9d1d9',
        'osc-bright': '#e6edf3',
        'osc-green': '#00e676',
        'osc-green-glow': 'rgba(0, 230, 118, 0.35)',
        'osc-cyan': '#00e5ff',
        'osc-cyan-glow': 'rgba(0, 229, 255, 0.3)',
        'osc-red': '#ff3d3d',
        'osc-red-glow': 'rgba(255, 61, 61, 0.4)',
        'osc-amber': '#ffb300',
        'osc-amber-glow': 'rgba(255, 179, 0, 0.3)',
      },
      fontFamily: {
        mono: ['"JetBrains Mono"', '"Space Mono"', 'ui-monospace', 'SFMono-Regular', 'monospace'],
        display: ['"JetBrains Mono"', '"Space Mono"', 'ui-monospace', 'monospace'],
      },
      boxShadow: {
        'green-glow': '0 0 20px rgba(0, 230, 118, 0.4), 0 0 40px rgba(0, 230, 118, 0.15)',
        'cyan-glow': '0 0 20px rgba(0, 229, 255, 0.35)',
        'red-glow': '0 0 15px rgba(255, 61, 61, 0.5)',
        'card': '0 1px 3px rgba(0,0,0,0.3), 0 8px 24px rgba(0,0,0,0.2)',
      },
      animation: {
        'scan-line': 'scanLine 4s linear infinite',
        'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        'glow': 'glow 2s ease-in-out infinite alternate',
      },
      keyframes: {
        scanLine: {
          '0%': { transform: 'translateY(-100%)' },
          '100%': { transform: 'translateY(100%)' },
        },
        glow: {
          '0%': { boxShadow: '0 0 5px rgba(0, 230, 118, 0.3)' },
          '100%': { boxShadow: '0 0 20px rgba(0, 230, 118, 0.7), 0 0 40px rgba(0, 230, 118, 0.3)' },
        },
      },
    },
  },
  plugins: [],
};
