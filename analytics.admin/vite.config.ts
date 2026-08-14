import { readdirSync, readFileSync } from 'node:fs';
import path from 'node:path';
import { tanstackRouter } from '@tanstack/router-plugin/vite';
import react from '@vitejs/plugin-react';
import { defineConfig, type Plugin } from 'vite';

const localesDir = path.resolve(import.meta.dirname, 'src/shared/i18n/locales');
const shimDir = path.resolve(import.meta.dirname, 'src/shared/lib/react-shims');

const REACT_SHIM_SPECIFIERS = {
  react: path.join(shimDir, 'react.ts'),
  'react-dom': path.join(shimDir, 'react-dom.ts'),
  'react-dom/client': path.join(shimDir, 'react-dom-client.ts'),
  'react/jsx-runtime': path.join(shimDir, 'jsx-runtime.ts'),
  'react/jsx-dev-runtime': path.join(shimDir, 'jsx-dev-runtime.ts'),
} as const;

function normalizeBasePath(basePath: string) {
  if (!basePath || basePath === '/') {
    return '/';
  }

  const withLeadingSlash = basePath.startsWith('/') ? basePath : `/${basePath}`;
  return withLeadingSlash.endsWith('/') ? withLeadingSlash : `${withLeadingSlash}/`;
}

function withBasePath(pathname: string, base: string) {
  const normalizedPath = pathname.startsWith('/') ? pathname.slice(1) : pathname;
  return `${base}${normalizedPath}`;
}

function localesJsonPlugin(base: string): Plugin {
  const basePrefix = base === '/' ? '' : base.replace(/\/$/, '');
  const localesPattern = new RegExp(`^${basePrefix}/locales/([a-z]+)\\.json(?:\\?.*)?$`);

  return {
    name: 'locales-json',
    configureServer(server) {
      server.middlewares.use((req, res, next) => {
        const match = req.url?.match(localesPattern);
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

function reactImportMapDevPlugin(base: string): Plugin {
  const devShimImports = {
    react: withBasePath('src/shared/lib/react-shims/react.ts', base),
    'react-dom': withBasePath('src/shared/lib/react-shims/react-dom.ts', base),
    'react-dom/client': withBasePath('src/shared/lib/react-shims/react-dom-client.ts', base),
    'react/jsx-runtime': withBasePath('src/shared/lib/react-shims/jsx-runtime.ts', base),
    'react/jsx-dev-runtime': withBasePath('src/shared/lib/react-shims/jsx-dev-runtime.ts', base),
  };

  return {
    name: 'react-importmap-dev',
    transformIndexHtml: {
      order: 'pre',
      handler(html) {
        return html.replace(
          /<script type="importmap">[\s\S]*?<\/script>/,
          `<script type="importmap">${JSON.stringify({ imports: devShimImports })}</script>`,
        );
      },
    },
  };
}

function reactImportMapBuildPlugin(base: string): Plugin {
  return {
    name: 'react-importmap-build',
    transformIndexHtml: {
      order: 'post',
      handler(html, ctx) {
        if (!ctx.bundle) {
          return html;
        }

        const imports: Record<string, string> = {};

        for (const [specifier, filePath] of Object.entries(REACT_SHIM_SPECIFIERS)) {
          const chunk = Object.values(ctx.bundle).find(
            (item) => item.type === 'chunk' && item.facadeModuleId === filePath,
          );

          if (chunk?.type === 'chunk') {
            imports[specifier] = withBasePath(chunk.fileName, base);
          }
        }

        return html.replace(
          /<script type="importmap">[\s\S]*?<\/script>/,
          `<script type="importmap">${JSON.stringify({ imports })}</script>`,
        );
      },
    },
  };
}

export default defineConfig(() => {
  const base = normalizeBasePath(process.env.VITE_BASE_PATH ?? '/');

  return {
    base,
    plugins: [
      tanstackRouter({
        target: 'react',
        autoCodeSplitting: true,
        routesDirectory: './src/app/routes',
        generatedRouteTree: './src/app/routeTree.gen.ts',
      }),
      react(),
      localesJsonPlugin(base),
      reactImportMapDevPlugin(base),
      reactImportMapBuildPlugin(base),
    ],
    resolve: {
      alias: {
        '@': path.resolve(import.meta.dirname, './src'),
      },
    },
    build: {
      rollupOptions: {
        input: {
          main: path.resolve(import.meta.dirname, 'index.html'),
          'shim-react': REACT_SHIM_SPECIFIERS.react,
          'shim-react-dom': REACT_SHIM_SPECIFIERS['react-dom'],
          'shim-react-dom-client': REACT_SHIM_SPECIFIERS['react-dom/client'],
          'shim-jsx-runtime': REACT_SHIM_SPECIFIERS['react/jsx-runtime'],
          'shim-jsx-dev-runtime': REACT_SHIM_SPECIFIERS['react/jsx-dev-runtime'],
        },
      },
    },
    server: {
      proxy: {
        '/api': { target: 'https://recaps.hakolr.dev', changeOrigin: true },
      },
    },
  };
});
