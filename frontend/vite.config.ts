import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    host: true,
    allowedHosts: ["short.irvineafri.com"],
    hmr: {
      overlay: false, // Disable error overlay that might cause refreshes
    },
    watch: {
      // Ignore node_modules and other files that shouldn't trigger reloads
      ignored: ['**/node_modules/**', '**/.git/**'],
    },
  },
  define: {
    // Define environment variables with defaults
    'import.meta.env.VITE_API_URL': JSON.stringify(process.env.VITE_API_URL || 'https://s.iafri.com'),
  },
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['react', 'react-dom'],
          router: ['react-router-dom'],
        },
      },
    },
  },
}) 