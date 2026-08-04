import type { ResolvedTheme, ThemeTokens } from '../types/theme';

const DEFAULT_THEME: ResolvedTheme = {
  colors: {
    bg: '#0B1F33',
    fg: '#FFFFFF',
    muted: 'rgba(255,255,255,0.72)',
    accent: '#00AAFF',
    accentSoft: 'rgba(0,170,255,0.16)',
    surface: 'rgba(255,255,255,0.06)',
    callout: 'rgba(0,170,255,0.12)',
  },
  fonts: {
    display: '"Manrope", system-ui, sans-serif',
    body: '"Manrope", system-ui, sans-serif',
  },
  radii: {
    button: 16,
    card: 24,
  },
  space: {
    scenePadding: 24,
  },
  assets: {},
  cssVars: {},
};

function toCssVars(theme: Omit<ResolvedTheme, 'cssVars'>): Record<string, string> {
  return {
    '--recap-bg': theme.colors.bg,
    '--recap-fg': theme.colors.fg,
    '--recap-muted': theme.colors.muted,
    '--recap-accent': theme.colors.accent,
    '--recap-accent-soft': theme.colors.accentSoft,
    '--recap-surface': theme.colors.surface,
    '--recap-callout': theme.colors.callout,
    '--recap-font-display': theme.fonts.display,
    '--recap-font-body': theme.fonts.body,
    '--recap-radius-button': `${theme.radii.button}px`,
    '--recap-radius-card': `${theme.radii.card}px`,
    '--recap-scene-padding': `${theme.space.scenePadding}px`,
    ...(theme.assets.background ? { '--recap-bg-image': `url(${theme.assets.background})` } : {}),
  };
}

export function createTheme(tokens: ThemeTokens = {}): ResolvedTheme {
  const resolved: Omit<ResolvedTheme, 'cssVars'> = {
    colors: {
      ...DEFAULT_THEME.colors,
      ...tokens.colors,
    },
    fonts: {
      ...DEFAULT_THEME.fonts,
      ...tokens.fonts,
    },
    radii: {
      ...DEFAULT_THEME.radii,
      ...tokens.radii,
    },
    space: {
      ...DEFAULT_THEME.space,
      ...tokens.space,
    },
    assets: {
      ...DEFAULT_THEME.assets,
      ...tokens.assets,
    },
  };

  return {
    ...resolved,
    cssVars: toCssVars(resolved),
  };
}
