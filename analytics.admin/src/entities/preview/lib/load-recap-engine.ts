import type { ComponentType } from 'react';
import * as React from 'react';
import * as JsxDevRuntime from 'react/jsx-dev-runtime';
import * as JsxRuntime from 'react/jsx-runtime';
import * as ReactDOM from 'react-dom';
import * as ReactDOMClient from 'react-dom/client';

export type RecapPalette = {
  bg: string;
  fg: string;
  muted: string;
  accent: string;
  accentSoft: string;
  surface: string;
  callout: string;
};

export type PreparedRecap = {
  data: unknown;
  scenes: unknown[];
  badges?: unknown[];
  locale?: string;
};

export type RecapEngineModule = {
  Recap: ComponentType<Record<string, unknown>>;
  prepareRecap: (payload: unknown) => PreparedRecap;
  createTheme: (tokens?: { colors?: RecapPalette }) => unknown;
};

type RecapEngineImporter = (url: string) => Promise<unknown>;

const cache = new Map<string, RecapEngineModule>();
const materializedBlobs = new Map<string, string>();
const inflightMaterialize = new Map<string, Promise<string>>();
const STYLES_ID = 'recap-engine-styles';

type ReactShimUrls = {
  react: string;
  'react-dom': string;
  'react-dom/client': string;
  'react/jsx-runtime': string;
  'react/jsx-dev-runtime': string;
};

let shimUrls: ReactShimUrls | null = null;

export function getRecapEngineImportUrl(version: string) {
  return `https://esm.sh/@recap-engine/react@${version}?external=react,react-dom`;
}

export function getRecapEngineStylesUrl(version: string) {
  return `https://esm.sh/@recap-engine/react@${version}/styles.css`;
}

function isIdentifier(name: string) {
  return /^[A-Za-z_$][\w$]*$/.test(name);
}

function namespaceToBlobUrl(module: object) {
  type ShimStore = Map<string, object>;
  const globalWithShims = globalThis as typeof globalThis & { __recapEngineShims?: ShimStore };
  const shims = globalWithShims.__recapEngineShims ?? new Map<string, object>();
  globalWithShims.__recapEngineShims = shims;

  const key = `shim:${shims.size}:${Math.random().toString(36).slice(2)}`;
  shims.set(key, module);

  const names = Object.keys(module).filter((name) => name !== 'default' && isIdentifier(name));
  const source = `
const m = globalThis.__recapEngineShims.get(${JSON.stringify(key)});
${names.map((name) => `export const ${name} = m[${JSON.stringify(name)}];`).join('\n')}
export default m.default ?? m;
`;

  return URL.createObjectURL(new Blob([source], { type: 'text/javascript' }));
}

function getReactShimUrls() {
  if (!shimUrls) {
    shimUrls = {
      react: namespaceToBlobUrl(React),
      'react-dom': namespaceToBlobUrl(ReactDOM),
      'react-dom/client': namespaceToBlobUrl(ReactDOMClient),
      'react/jsx-runtime': namespaceToBlobUrl(JsxRuntime),
      'react/jsx-dev-runtime': namespaceToBlobUrl(JsxDevRuntime),
    };
  }

  return shimUrls;
}

export function rewriteEsmShSource(source: string, shims: Record<string, string>) {
  let next = source;
  const specifiers = Object.keys(shims).sort((left, right) => right.length - left.length);

  for (const specifier of specifiers) {
    const blobUrl = shims[specifier];
    if (!blobUrl) {
      continue;
    }

    next = next.replaceAll(`"${specifier}"`, `"${blobUrl}"`);
    next = next.replaceAll(`'${specifier}'`, `'${blobUrl}'`);
  }

  next = next.replace(
    /\bfrom\s*(["'])(\/[^"']+)\1/g,
    (_match, quote: string, pathname: string) => `from ${quote}https://esm.sh${pathname}${quote}`,
  );
  next = next.replace(
    /\bimport\s*(["'])(\/[^"']+)\1/g,
    (_match, quote: string, pathname: string) => `import ${quote}https://esm.sh${pathname}${quote}`,
  );
  next = next.replace(
    /\bexport\s*\*\s*from\s*(["'])(\/[^"']+)\1/g,
    (_match, quote: string, pathname: string) =>
      `export * from ${quote}https://esm.sh${pathname}${quote}`,
  );

  return next;
}

export function collectEsmShUrls(source: string) {
  const matches = source.match(/https:\/\/esm\.sh\/[^"'\s]+/g) ?? [];
  return [...new Set(matches)].sort((left, right) => right.length - left.length);
}

async function materializeEsmShModule(url: string): Promise<string> {
  const cached = materializedBlobs.get(url);
  if (cached) {
    return cached;
  }

  const pending = inflightMaterialize.get(url);
  if (pending) {
    return pending;
  }

  const job = (async () => {
    const response = await fetch(url);
    if (!response.ok) {
      throw new Error(`Failed to load recap engine from ${url}: ${response.status}`);
    }

    let source = rewriteEsmShSource(await response.text(), getReactShimUrls());
    const deps = collectEsmShUrls(source).filter(
      (dep) => !dep.endsWith('.css') && !dep.endsWith('.map'),
    );

    for (const dep of deps) {
      const depBlob = await materializeEsmShModule(dep);
      source = source.replaceAll(dep, depBlob);
    }

    const blobUrl = URL.createObjectURL(new Blob([source], { type: 'text/javascript' }));
    materializedBlobs.set(url, blobUrl);
    return blobUrl;
  })();

  inflightMaterialize.set(url, job);
  try {
    return await job;
  } finally {
    inflightMaterialize.delete(url);
  }
}

async function importMaterializedEsmSh(url: string) {
  const blobUrl = await materializeEsmShModule(url);
  return import(/* @vite-ignore */ blobUrl);
}

export function applyRecapEngineStyles(version: string) {
  if (typeof document === 'undefined') {
    return;
  }

  let link = document.getElementById(STYLES_ID) as HTMLLinkElement | null;
  if (!link) {
    link = document.createElement('link');
    link.id = STYLES_ID;
    link.rel = 'stylesheet';
    document.head.appendChild(link);
  }

  const href = getRecapEngineStylesUrl(version);
  link.setAttribute('data-href', href);
  if (import.meta.env.MODE !== 'test') {
    link.href = href;
  }
}

function isRecapEngineModule(value: unknown): value is RecapEngineModule {
  if (!value || typeof value !== 'object') {
    return false;
  }

  const module = value as RecapEngineModule;
  const recapType = typeof module.Recap;

  return (
    (recapType === 'function' || recapType === 'object') &&
    Boolean(module.Recap) &&
    typeof module.prepareRecap === 'function' &&
    typeof module.createTheme === 'function'
  );
}

export async function loadRecapEngine(
  version: string,
  importer: RecapEngineImporter = importMaterializedEsmSh,
): Promise<RecapEngineModule> {
  const cached = cache.get(version);
  if (cached) {
    applyRecapEngineStyles(version);
    return cached;
  }

  const imported = await importer(getRecapEngineImportUrl(version));
  if (!isRecapEngineModule(imported)) {
    throw new Error('Invalid recap engine module');
  }

  cache.set(version, imported);
  applyRecapEngineStyles(version);
  return imported;
}

export function resetRecapEngineCache() {
  cache.clear();
  for (const blobUrl of materializedBlobs.values()) {
    URL.revokeObjectURL(blobUrl);
  }
  materializedBlobs.clear();
}
