import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/system': {
        target: 'http://localhost:8070',
        changeOrigin: true,
      },
      '/api': {
        target: 'http://localhost:8070',
        changeOrigin: true,
      }
    }
  }
})
