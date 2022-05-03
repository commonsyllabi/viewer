import { defineConfig } from 'vite'
import { fileURLToPath } from 'url'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  root: './src',
  build: {
    outDir: '../public',
    rollupOptions: {
      input: {
        home: fileURLToPath(new URL('./src/index.html', import.meta.url)),
        cartridge: fileURLToPath(new URL('./src/cartridge.html', import.meta.url)),
      }
    }
  }
})
