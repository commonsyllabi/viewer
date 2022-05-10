import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import { resolve } from 'path'
import { fileURLToPath } from 'url'

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
    outDir: '../public',
    rollupOptions: {
      input: {
        home: fileURLToPath(new URL('./src/index.html', import.meta.url)),
        cartridge: fileURLToPath(new URL('./src/cartridge.html', import.meta.url)),
        syllabus: fileURLToPath(new URL('./src/syllabus.html', import.meta.url)),
      },
      output: {
        entryFileNames: `assets/[name].js`,
        chunkFileNames: `assets/[name].js`,
        assetFileNames: `assets/[name].[ext]`
      }
    }
  }
})
