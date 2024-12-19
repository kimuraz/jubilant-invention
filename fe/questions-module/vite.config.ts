import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  build: {
    lib: {
      entry: './src/ce.ts',
      name: 'questions',
      fileName: 'questions',
      formats: ['es']
    },
    emptyOutDir: false,
    outDir: '../static-server/fe/',
  },
  define: {
    'process.env': process.env
  }
})
