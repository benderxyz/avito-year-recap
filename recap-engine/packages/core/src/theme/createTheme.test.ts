import { describe, expect, it } from '@jest/globals';

import { createTheme } from './createTheme';

describe('createTheme', () => {
  it('returns all default tokens and matching CSS variables', () => {
    const theme = createTheme();

    expect(theme.colors).toEqual({
      bg: '#0B1F33',
      fg: '#FFFFFF',
      muted: 'rgba(255,255,255,0.72)',
      accent: '#00AAFF',
      accentSoft: 'rgba(0,170,255,0.16)',
      surface: 'rgba(255,255,255,0.06)',
      callout: 'rgba(0,170,255,0.12)',
    });
    expect(theme.fonts).toEqual({
      display: '"Manrope", system-ui, sans-serif',
      body: '"Manrope", system-ui, sans-serif',
    });
    expect(theme.radii).toEqual({ button: 16, card: 24 });
    expect(theme.space).toEqual({ scenePadding: 24 });
    expect(theme.assets).toEqual({});
    expect(theme.cssVars).toEqual({
      '--recap-bg': '#0B1F33',
      '--recap-fg': '#FFFFFF',
      '--recap-muted': 'rgba(255,255,255,0.72)',
      '--recap-accent': '#00AAFF',
      '--recap-accent-soft': 'rgba(0,170,255,0.16)',
      '--recap-surface': 'rgba(255,255,255,0.06)',
      '--recap-callout': 'rgba(0,170,255,0.12)',
      '--recap-font-display': '"Manrope", system-ui, sans-serif',
      '--recap-font-body': '"Manrope", system-ui, sans-serif',
      '--recap-radius-button': '16px',
      '--recap-radius-card': '24px',
      '--recap-scene-padding': '24px',
    });
  });

  it('deep-merges partial overrides and derives CSS variables including background image', () => {
    const theme = createTheme({
      colors: { bg: '#101010', accent: '#ff0' },
      fonts: { display: 'Display' },
      radii: { button: 0 },
      space: { scenePadding: 40 },
      assets: { background: '/background image.png' },
    });

    expect(theme.colors).toMatchObject({ bg: '#101010', fg: '#FFFFFF', accent: '#ff0' });
    expect(theme.fonts).toEqual({ display: 'Display', body: '"Manrope", system-ui, sans-serif' });
    expect(theme.radii).toEqual({ button: 0, card: 24 });
    expect(theme.space.scenePadding).toBe(40);
    expect(theme.assets.background).toBe('/background image.png');
    expect(theme.cssVars).toMatchObject({
      '--recap-bg': '#101010',
      '--recap-accent': '#ff0',
      '--recap-font-display': 'Display',
      '--recap-radius-button': '0px',
      '--recap-scene-padding': '40px',
      '--recap-bg-image': 'url(/background image.png)',
    });
  });
});
