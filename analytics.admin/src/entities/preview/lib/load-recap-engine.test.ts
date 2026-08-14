import { afterEach, describe, expect, it, vi } from 'vitest';
import {
  collectEsmShUrls,
  getRecapEngineImportUrl,
  getRecapEngineStylesUrl,
  loadRecapEngine,
  resetRecapEngineCache,
  rewriteEsmShSource,
} from './load-recap-engine';

const engineModule = {
  Recap: () => null,
  prepareRecap: () => ({ data: {}, scenes: [] }),
  createTheme: () => ({ cssVars: {} }),
};

describe('loadRecapEngine', () => {
  afterEach(() => {
    resetRecapEngineCache();
    document.getElementById('recap-engine-styles')?.remove();
  });

  it('builds esm.sh URLs', () => {
    expect(getRecapEngineImportUrl('2.0.1')).toBe(
      'https://esm.sh/@recap-engine/react@2.0.1?external=react,react-dom',
    );
    expect(getRecapEngineStylesUrl('2.0.1')).toBe(
      'https://esm.sh/@recap-engine/react@2.0.1/styles.css',
    );
  });

  it('rewrites bare react specifiers and esm.sh absolute paths', () => {
    const source = `import "/@recap-engine/core@1.3.0/es2022/core.mjs";
import{useContext}from"react";
import{jsx}from"react/jsx-runtime";
export * from "/@recap-engine/react@2.0.1/es2022/react.mjs";`;

    const rewritten = rewriteEsmShSource(source, {
      react: 'blob:react',
      'react/jsx-runtime': 'blob:jsx',
    });

    expect(rewritten).toContain('from"blob:react"');
    expect(rewritten).toContain('from"blob:jsx"');
    expect(rewritten).toContain('import "https://esm.sh/@recap-engine/core@1.3.0/es2022/core.mjs"');
    expect(rewritten).toContain(
      'export * from "https://esm.sh/@recap-engine/react@2.0.1/es2022/react.mjs"',
    );
    expect(collectEsmShUrls(rewritten)).toEqual([
      'https://esm.sh/@recap-engine/react@2.0.1/es2022/react.mjs',
      'https://esm.sh/@recap-engine/core@1.3.0/es2022/core.mjs',
    ]);
  });

  it('imports, caches, and injects styles', async () => {
    const importer = vi.fn(async () => engineModule);

    const first = await loadRecapEngine('2.0.1', importer);
    const second = await loadRecapEngine('2.0.1', importer);

    expect(first).toBe(engineModule);
    expect(second).toBe(engineModule);
    expect(importer).toHaveBeenCalledOnce();
    expect(importer).toHaveBeenCalledWith(
      'https://esm.sh/@recap-engine/react@2.0.1?external=react,react-dom',
    );
    expect(document.getElementById('recap-engine-styles')).toHaveAttribute(
      'data-href',
      'https://esm.sh/@recap-engine/react@2.0.1/styles.css',
    );
  });

  it('rejects a module without Recap exports', async () => {
    await expect(loadRecapEngine('2.0.1', async () => ({}))).rejects.toThrow(
      'Invalid recap engine module',
    );
  });
});
