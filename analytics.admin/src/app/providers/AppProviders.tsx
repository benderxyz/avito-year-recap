import type { ReactNode } from 'react';
import { MantineProvider, mantineTheme } from '../../shared/ui/theme/mantine-theme';
import I18nProvider from '@/shared/i18n/I18nProvider';

type AppProvidersProps = {
  children?: ReactNode;
};

export default function AppProviders({ children }: AppProvidersProps) {
  return (
    <MantineProvider theme={mantineTheme} defaultColorScheme="light">
      <I18nProvider>{children}</I18nProvider>
    </MantineProvider>
  );
}
