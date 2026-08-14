import { describe, expect, it } from 'vitest';
import { compareSemverDesc, getCoreDependency, parseReactPackument } from './parse-packument';

describe('parseReactPackument', () => {
  const packument = {
    'dist-tags': { latest: '2.0.1' },
    versions: {
      '1.0.0': { dependencies: { '@recap-engine/core': '1.0.0' } },
      '2.0.1': { dependencies: { '@recap-engine/core': '^1.3.0' } },
      '2.0.0': { dependencies: {} },
      '2.1.0-beta.1': { dependencies: { '@recap-engine/core': '1.4.0' } },
    },
  };

  it('sorts versions descending and reads the latest tag', () => {
    const parsed = parseReactPackument(packument);

    expect(parsed.latest).toBe('2.0.1');
    expect(parsed.versions).toEqual(['2.1.0-beta.1', '2.0.1', '2.0.0', '1.0.0']);
  });

  it('reads the core dependency specifier', () => {
    const parsed = parseReactPackument(packument);

    expect(getCoreDependency(parsed, '2.0.1')).toBe('^1.3.0');
    expect(getCoreDependency(parsed, '2.0.0')).toBeNull();
    expect(getCoreDependency(parsed, 'missing')).toBeNull();
  });
});

describe('compareSemverDesc', () => {
  it('places a stable release above a prerelease of the same patch', () => {
    expect(compareSemverDesc('2.0.1', '2.0.1-beta.1')).toBeLessThan(0);
  });
});
