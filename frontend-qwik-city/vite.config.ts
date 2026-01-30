import { defineConfig } from "vite";
import { qwikVite } from "@builder.io/qwik/optimizer";
import { qwikCity } from "@builder.io/qwik-city/vite";
import tsconfigPaths from "vite-tsconfig-paths";

export default defineConfig(() => {
  return {
    plugins: [qwikCity(), qwikVite(), tsconfigPaths()],
    server: {
      port: 3000,
      host: true,
      proxy: {
        "/api": {
          target: process.env.API_URL || "http://localhost:3001",
          changeOrigin: true,
        },
        "/uploads": {
          target: process.env.API_URL || "http://localhost:3001",
          changeOrigin: true,
        },
      },
    },
    preview: {
      port: 3000,
      host: true,
    },
  };
});
