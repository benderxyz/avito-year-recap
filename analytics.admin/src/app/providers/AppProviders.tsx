import { useState, type ReactNode } from 'react';

import { QueryClientProvider } from '@tanstack/react-query';
import I18nProvider from '@/shared/i18n/I18nProvider';
import { createQueryClient } from '@/shared/api/query-client';

import { MantineProvider, mantineTheme } from '@/shared/ui';

type AppProvidersProps = {
  children?: ReactNode;
};

export default function AppProviders({ children }: AppProvidersProps) {
  const [queryClient] = useState(createQueryClient);

  return (
    <QueryClientProvider client={queryClient}>
      <MantineProvider theme={mantineTheme} defaultColorScheme="light">
        <I18nProvider>{children}</I18nProvider>
      </MantineProvider>
    </QueryClientProvider>
  );
}
