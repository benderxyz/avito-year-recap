import type { NpmPackument } from '../api/npm';

export type ParsedReactPackument = {
  versions: string[];
  latest: string | null;
  coreByReactVersion: Record<string, string | null>;
};

function parseSemver(version: string) {
  const [core, pre] = version.split('-');
  const [major, minor, patch] = (core ?? '')
    .split('.')
    .map((part) => Number.parseInt(part, 10) || 0);

  return { major, minor, patch, pre: pre ?? '' };
}

export function compareSemverDesc(a: string, b: string) {
  const left = parseSemver(a);
  const right = parseSemver(b);

  if (left.major !== right.major) {
    return right.major - left.major;
  }
  if (left.minor !== right.minor) {
    return right.minor - left.minor;
  }
  if (left.patch !== right.patch) {
    return right.patch - left.patch;
  }
  if (!left.pre && right.pre) {
    return -1;
  }
  if (left.pre && !right.pre) {
    return 1;
  }

  return right.pre.localeCompare(left.pre);
}

export function parseReactPackument(packument: NpmPackument): ParsedReactPackument {
  const versions = Object.keys(packument.versions ?? {}).sort(compareSemverDesc);
  const latest = packument['dist-tags']?.latest ?? versions[0] ?? null;
  const coreByReactVersion: Record<string, string | null> = {};

  for (const version of versions) {
    coreByReactVersion[version] =
      packument.versions?.[version]?.dependencies?.['@recap-engine/core'] ?? null;
  }

  return { versions, latest, coreByReactVersion };
}

export function getCoreDependency(parsed: ParsedReactPackument, reactVersion: string) {
  return parsed.coreByReactVersion[reactVersion] ?? null;
}
