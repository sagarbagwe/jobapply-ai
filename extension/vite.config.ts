import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import { resolve } from 'node:path';

export default defineConfig({
  plugins: [react()],
  build: {
    outDir: 'dist',
    rollupOptions: {
      input: { popup: resolve(import.meta.dirname, 'index.html'), background: resolve(import.meta.dirname, 'src/background/index.ts'), content: resolve(import.meta.dirname, 'src/content/index.ts') },
      output: { entryFileNames: '[name].js' }
    }
  }
});
