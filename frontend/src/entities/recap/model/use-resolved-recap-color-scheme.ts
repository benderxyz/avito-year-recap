import { useMantineColorScheme } from '@mantine/core';
import { useColorScheme } from '@mantine/hooks';
import type { RecapColorScheme } from '../lib/recap-theme';
import { resolveRecapColorScheme } from '../lib/resolve-recap-color-scheme';

export function useResolvedRecapColorScheme(): RecapColorScheme {
  const { colorScheme } = useMantineColorScheme();
  const systemColorScheme = useColorScheme();

  return resolveRecapColorScheme(colorScheme, systemColorScheme);
}
