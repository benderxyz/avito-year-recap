import { createTheme } from '@recap-engine/core';
import { ERecapEventType, ESceneActionType, type RecapEvent } from '@recap-engine/react';

export type RecapColorScheme = 'light' | 'dark';

const RECAP_PALETTES: Record<
  RecapColorScheme,
  {
    bg: string;
    fg: string;
    muted: string;
    accent: string;
    accentSoft: string;
    surface: string;
    callout: string;
  }
> = {
  light: {
    bg: '#ffffff',
    fg: '#1a1a1a',
    muted: '#5c5c5c',
    accent: '#00aaff',
    accentSoft: '#e8f7ff',
    surface: '#f2f1f0',
    callout: 'rgba(0, 170, 255, 0.12)',
  },
  dark: {
    bg: '#1f1f1f',
    fg: '#f5f5f5',
    muted: '#a3a3a3',
    accent: '#00aaff',
    accentSoft: '#0d2a3d',
    surface: 'rgba(255, 255, 255, 0.06)',
    callout: 'rgba(0, 170, 255, 0.16)',
  },
};

export function buildRecapTheme(colorScheme: RecapColorScheme = 'light') {
  const palette = RECAP_PALETTES[colorScheme];

  return createTheme({
    colors: palette,
  });
}

export function shouldCloseRecapOnEvent(event: RecapEvent): boolean {
  if (event.type === ERecapEventType.Complete) {
    return true;
  }

  return (
    event.type === ERecapEventType.Action &&
    event.action.type === ESceneActionType.Custom &&
    event.action.id === 'close-recap'
  );
}
