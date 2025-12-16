// uno.config.js
import {
  defineConfig,
  presetUno,
  presetAttributify,
  presetIcons,
} from "unocss";
import transformerDirectives from "@unocss/transformer-directives";
import transformerVariantGroup from "@unocss/transformer-variant-group";

export default defineConfig({
  presets: [
    presetUno(),
    presetAttributify(),
    presetIcons({
      scale: 1.2,
      extraProperties: {
        display: 'inline-block',
        'vertical-align': 'middle',
      },

      // ★ 最关键：手动指定 JSON loader（修复 Node 22 的 assert 限制）
      collections: {
        'material-symbols': () =>
          import('@iconify-json/material-symbols/icons.json', {
            with: { type: 'json' }
          }).then(i => i.default),
      },
    }),
  ],
  transformers: [transformerDirectives(), transformerVariantGroup()],

  shortcuts: [
    [
      "btn",
      "px-4 py-2 rounded-lg bg-primary text-white font-medium transition hover:bg-primary/80",
    ],
    ["glass", "backdrop-blur-xl bg-white/60 border border-white/30 shadow-md"],
    [
      "glass-dark",
      "backdrop-blur-xl bg-black/30 text-white border border-white/10 shadow",
    ],
    ["flex-center", "flex justify-center items-center"],
  ],

  theme: {
    colors: {
      primary: "#3B82F6",
      secondary: "#10B981",
      danger: "#EF4444",
      border: "#e0e0e0ff",
    },
  },

  safelist: [
    "bg-white/40",
    "bg-white/50",
    "bg-white/60",
    "text-primary",
    "text-secondary",
  ],
});
