import type { RecapColorScheme } from '../lib/recap-theme';

type MantineColorScheme = RecapColorScheme | 'auto';

export function resolveRecapColorScheme(
  colorScheme: MantineColorScheme,
  systemColorScheme: RecapColorScheme,
): RecapColorScheme {
  if (colorScheme === 'dark' || colorScheme === 'light') {
    return colorScheme;
  }

  return systemColorScheme;
}
