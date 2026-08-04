import path from 'node:path';
import { fileURLToPath } from 'node:url';
import react from '@vitejs/plugin-react';
import { defineConfig } from 'vite';

const root = path.dirname(fileURLToPath(import.meta.url));
const recapEngine = path.resolve(root, '../recap-engine/packages');

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@recap-engine/core': path.join(recapEngine, 'core/src/index.ts'),
      '@recap-engine/react/styles.css': path.join(recapEngine, 'react/src/styles.css'),
      '@recap-engine/react': path.join(recapEngine, 'react/src/index.ts'),
    },
    dedupe: ['react', 'react-dom'],
  },
  server: {
    port: 3000,
    fs: {
      allow: [root, path.resolve(root, '../recap-engine')],
    },
  },
});
