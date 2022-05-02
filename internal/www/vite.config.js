import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import { resolve } from 'path'

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
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src')
    }
  },
  build: {
    outDir: "../public",
  },
  devServer: {
    proxy: "http://localhost:2046/",
  },
});
