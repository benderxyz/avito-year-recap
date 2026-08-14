import { describe, expect, it } from 'vitest';
import { getRouterBasepath, normalizeBasePath, withBasePath } from './app-base';

describe('app-base', () => {
  it('normalizes base paths', () => {
    expect(normalizeBasePath('/')).toBe('/');
    expect(normalizeBasePath('/admin')).toBe('/admin/');
    expect(normalizeBasePath('/admin/')).toBe('/admin/');
  });

  it('derives router basepath', () => {
    expect(getRouterBasepath('/')).toBeUndefined();
    expect(getRouterBasepath('/admin/')).toBe('/admin');
  });

  it('prefixes asset paths', () => {
    expect(withBasePath('locales/en.json', '/admin/')).toBe('/admin/locales/en.json');
    expect(withBasePath('/assets/app.js', '/')).toBe('/assets/app.js');
  });
});
