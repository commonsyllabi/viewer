import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
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
  root: './src',
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
    },
  },
  build: {
    outDir: '../public',
    rollupOptions: {
      input: {
        about: fileURLToPath(new URL('./src/about.html', import.meta.url)),
        listing: fileURLToPath(new URL('./src/listing.html', import.meta.url)),
        cartridge: fileURLToPath(new URL('./src/index.html', import.meta.url)),
        syllabus: fileURLToPath(
          new URL('./src/syllabus.html', import.meta.url)
        ),
        magic_link: fileURLToPath(
          new URL('./src/magic_link.html', import.meta.url)
        ),
        error: fileURLToPath(new URL('./src/error.html', import.meta.url)),
      },
      output: {
        entryFileNames: `assets/[name].js`,
        chunkFileNames: `assets/[name].js`,
        assetFileNames: `assets/[name].[ext]`,
      },
    },
  },
})
