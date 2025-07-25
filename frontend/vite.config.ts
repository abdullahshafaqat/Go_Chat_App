import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    host: true,
    port: 8080,
    allowedHosts: ['localhost', '127.0.0.1', 'chat-app.local'], 
  },
  build: {
    outDir: 'dist',
    emptyOutDir: true,
  },
})
