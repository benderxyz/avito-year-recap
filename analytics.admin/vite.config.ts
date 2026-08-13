import { readFileSync, readdirSync } from 'node:fs';
import path from 'node:path';
import { tanstackRouter } from '@tanstack/router-plugin/vite';
import react from '@vitejs/plugin-react';
import { defineConfig, type Plugin } from 'vite';

const localesDir = path.resolve(import.meta.dirname, 'src/shared/i18n/locales');

function localesJsonPlugin(): Plugin {
  return {
    name: 'locales-json',
    configureServer(server) {
      server.middlewares.use((req, res, next) => {
        const match = req.url?.match(/^\/locales\/([a-z]+)\.json(?:\?.*)?$/);
        if (!match) {
          next();
          return;
        }

        const locale = match[1];
        const file = path.join(localesDir, `${locale}.json`);

        try {
          res.setHeader('Content-Type', 'application/json');
          res.end(readFileSync(file));
        } catch {
          next();
        }
      });
    },
    generateBundle() {
      for (const fileName of readdirSync(localesDir)) {
        if (!fileName.endsWith('.json')) {
          continue;
        }

        this.emitFile({
          type: 'asset',
          fileName: `locales/${fileName}`,
          source: readFileSync(path.join(localesDir, fileName)),
        });
      }
    },
  };
}

export default defineConfig({
  plugins: [
    tanstackRouter({
      target: 'react',
      autoCodeSplitting: true,
      routesDirectory: './src/app/routes',
      generatedRouteTree: './src/app/routeTree.gen.ts',
    }),
    react(),
    localesJsonPlugin(),
  ],
  resolve: {
    alias: {
      '@': path.resolve(import.meta.dirname, './src'),
    },
  },
  server: {
    proxy: {
      '/api': { target: 'https://recaps.hakolr.dev/api', changeOrigin: true },
    },
  },
});
