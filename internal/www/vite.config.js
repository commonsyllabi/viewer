import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  //   preprocessorOptions: {
  //     sass: {
  //       additionalData: '@import "@/css/global-vars.scss";'
  //     }
  //   }
  // },
  root: "./src",
  build: {
    outDir: "../public",
  },
  devServer: {
    proxy: "http://localhost:2046/",
  },
});
