export function normalizeBasePath(basePath: string) {
  if (!basePath || basePath === '/') {
    return '/';
  }

  const withLeadingSlash = basePath.startsWith('/') ? basePath : `/${basePath}`;
  return withLeadingSlash.endsWith('/') ? withLeadingSlash : `${withLeadingSlash}/`;
}

export function getRouterBasepath(baseUrl = import.meta.env.BASE_URL) {
  const normalized = normalizeBasePath(baseUrl);
  return normalized === '/' ? undefined : normalized.slice(0, -1);
}

export function withBasePath(path: string, baseUrl = import.meta.env.BASE_URL) {
  const normalized = normalizeBasePath(baseUrl);
  const normalizedPath = path.startsWith('/') ? path.slice(1) : path;
  return `${normalized}${normalizedPath}`;
}
