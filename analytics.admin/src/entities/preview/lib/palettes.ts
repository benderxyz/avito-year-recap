import type { RecapPalette } from './load-recap-engine';

export const PREVIEW_PALETTE_IDS = [
  'avitoLight',
  'avitoDark',
  'midnight',
  'sunset',
  'forest',
] as const;

export type PreviewPaletteId = (typeof PREVIEW_PALETTE_IDS)[number];

export const PREVIEW_PALETTES: Record<PreviewPaletteId, RecapPalette> = {
  avitoLight: {
    bg: '#ffffff',
    fg: '#1a1a1a',
    muted: '#5c5c5c',
    accent: '#00aaff',
    accentSoft: '#e8f7ff',
    surface: '#f2f1f0',
    callout: 'rgba(0, 170, 255, 0.12)',
  },
  avitoDark: {
    bg: '#1f1f1f',
    fg: '#f5f5f5',
    muted: '#a3a3a3',
    accent: '#00aaff',
    accentSoft: '#0d2a3d',
    surface: 'rgba(255, 255, 255, 0.06)',
    callout: 'rgba(0, 170, 255, 0.16)',
  },
  midnight: {
    bg: '#0B1F33',
    fg: '#FFFFFF',
    muted: 'rgba(255,255,255,0.72)',
    accent: '#00AAFF',
    accentSoft: 'rgba(0,170,255,0.16)',
    surface: 'rgba(255,255,255,0.06)',
    callout: 'rgba(0,170,255,0.12)',
  },
  sunset: {
    bg: '#1A0F0A',
    fg: '#FFF6EE',
    muted: 'rgba(255,246,238,0.72)',
    accent: '#FF6B35',
    accentSoft: 'rgba(255,107,53,0.18)',
    surface: 'rgba(255,255,255,0.08)',
    callout: 'rgba(255,107,53,0.14)',
  },
  forest: {
    bg: '#0F1F14',
    fg: '#E8F8EE',
    muted: 'rgba(232,248,238,0.72)',
    accent: '#3DDC97',
    accentSoft: 'rgba(61,220,151,0.16)',
    surface: 'rgba(255,255,255,0.06)',
    callout: 'rgba(61,220,151,0.12)',
  },
};
