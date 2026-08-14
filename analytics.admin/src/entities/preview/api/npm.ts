export const REACT_PACKUMENT_URL = 'https://registry.npmjs.org/@recap-engine/react';

export type NpmPackumentVersion = {
  dependencies?: Record<string, string>;
};

export type NpmPackument = {
  'dist-tags'?: {
    latest?: string;
  };
  versions?: Record<string, NpmPackumentVersion>;
};

export async function fetchReactPackument(): Promise<NpmPackument> {
  const response = await fetch(REACT_PACKUMENT_URL);
  if (!response.ok) {
    throw new Error(`Failed to load npm packument: ${response.status}`);
  }

  return response.json() as Promise<NpmPackument>;
}
