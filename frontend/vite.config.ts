import UnoCSS from "unocss/vite";
import { defineConfig, ProxyOptions } from "vite";
import vue from "@vitejs/plugin-vue";

const backend: ProxyOptions = {
  target: "http://localhost:3668",
  changeOrigin: true,
};
const base = import.meta.env.CONSEE_BASE || "";

// https://vite.dev/config/
export default defineConfig({
  base,
  plugins: [vue(), UnoCSS()],
  envPrefix: "CONSEE_",
  server: {
    proxy: {
      "^/api/v0/.*": backend,
    },
  },
});
